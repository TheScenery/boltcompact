package boltcompact

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

type CompactOption struct {
	TxMaxSize int64
	Debug     bool
}

// walkFunc is the type of the function called for keys (buckets and "normal"
// values) discovered by Walk. keys is the list of keys to descend to the bucket
// owning the discovered key/value pair k/v.
type walkFunc func(keys [][]byte, k, v []byte, seq uint64) error

func walkBucket(b *bolt.Bucket, keypath [][]byte, k, v []byte, seq uint64, fn walkFunc) error {
	// Execute callback.
	if err := fn(keypath, k, v, seq); err != nil {
		return err
	}

	// If this is not a bucket then stop.
	if v != nil {
		return nil
	}

	// Iterate over each child key/value.
	keypath = append(keypath, k)
	return b.ForEach(func(k, v []byte) error {
		if v == nil {
			bkt := b.Bucket(k)
			return walkBucket(bkt, keypath, k, nil, bkt.Sequence(), fn)
		}
		return walkBucket(b, keypath, k, v, b.Sequence(), fn)
	})
}

// walk walks recursively the bolt database db, calling walkFn for each key it finds.
func walk(db *bolt.DB, walkFn walkFunc) error {
	return db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			return walkBucket(b, nil, name, nil, b.Sequence(), walkFn)
		})
	})
}

func compactImp(dst, src *bolt.DB, opt *CompactOption) error {
	// commit regularly, or we'll run out of memory for large datasets if using one transaction.
	var size int64
	tx, err := dst.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := walk(src, func(keys [][]byte, k, v []byte, seq uint64) error {
		// On each key/value, check if we have exceeded tx size.
		sz := int64(len(k) + len(v))
		if size+sz > opt.TxMaxSize {
			// Commit previous transaction.
			if err := tx.Commit(); err != nil {
				return err
			}

			// Start new transaction.
			tx, err = dst.Begin(true)
			if err != nil {
				return err
			}
			size = 0
		}
		size += sz

		// Create bucket on the root transaction if this is the first level.
		nk := len(keys)
		if nk == 0 {
			bkt, err := tx.CreateBucket(k)
			if err != nil {
				return err
			}
			if err := bkt.SetSequence(seq); err != nil {
				return err
			}
			return nil
		}

		// Create buckets on subsequent levels, if necessary.
		b := tx.Bucket(keys[0])
		if nk > 1 {
			for _, k := range keys[1:] {
				b = b.Bucket(k)
			}
		}

		// If there is no value then this is a bucket call.
		if v == nil {
			bkt, err := b.CreateBucket(k)
			if err != nil {
				return err
			}
			if err := bkt.SetSequence(seq); err != nil {
				return err
			}
			return nil
		}

		// Otherwise treat it as a key/value pair.
		return b.Put(k, v)
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func Compact(src, dst string, opt *CompactOption) error {
	// Ensure source file exists.
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if dst == "" {
		return fmt.Errorf("dest path can not be empty")
	}

	// Open source database.
	srcDB, err := bolt.Open(src, 0444, nil)
	if err != nil {
		return err
	}
	defer srcDB.Close()

	// Open destination database.
	dstDB, err := bolt.Open(dst, srcFileInfo.Mode(), nil)
	if err != nil {
		return err
	}
	defer dstDB.Close()

	if opt == nil {
		opt = &CompactOption{
			TxMaxSize: 65536,
		}
	}

	if err := compactImp(dstDB, srcDB, opt); err != nil {
		return err
	}

	// Report stats on new size.
	dstFileInfo, err := os.Stat(dst)
	if err != nil {
		return err
	} else if dstFileInfo.Size() == 0 {
		return fmt.Errorf("zero db size")
	}
	if opt.Debug {
		initialSize := srcFileInfo.Size()
		compactedSize := dstFileInfo.Size()
		fmt.Printf("%d -> %d bytes (gain=%.2fx)\n", initialSize, compactedSize, float64(initialSize)/float64(compactedSize))
	}

	return nil
}
