package tar

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func Extract(reader io.Reader, target string, blacklist []string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("error == %s",err )
			return err
		}

		filename := header.Name
		filename = path.Join(target, header.Name)

		skip := false
		for _, prefix := range blacklist {
			if strings.HasPrefix(header.Name, prefix) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			log.Printf("1 copying %s",filename)
			if err = os.MkdirAll(filename, os.FileMode(header.Mode)); err != nil {
				log.Printf("error == %s",err )
				return err
			}
		case tar.TypeReg:
			log.Printf("2 copying %s",filename)
			if _, err := os.Stat(filename); err == nil {
				if err := os.Remove(filename); err != nil {
					log.Printf("error == %s",err )
					return err
				}
			}
			writer, err := os.Create(filename)
			if err != nil {
				log.Printf("error == %s",err )
				return err
			}
			io.Copy(writer, tarReader)
			if err = os.Chmod(filename, header.FileInfo().Mode()); err != nil {
				log.Printf("error == %s",err )
				return err
			}
			writer.Close()
		case tar.TypeLink:
			header.Linkname = "/" + header.Linkname
			log.Printf("3 copying actual = %s link = %s",header.Linkname, filename)
			if _, err := os.Stat(filename); err == nil {
				if err := os.Remove(filename); err != nil {
					log.Printf("error == %s",err )
					
				}
			}
			if err := os.Symlink(header.Linkname, filename); err != nil {
				log.Printf("error == %s",err )
			}
		case tar.TypeSymlink:
			log.Printf("4 copying actual = %s link = %s",header.Linkname, filename)
			if _, err := os.Stat(filename); err == nil {
				if err := os.Remove(filename); err != nil {
					log.Printf("error == %s",err )
					
				}
			}
			if err := os.Symlink(header.Linkname, filename); err != nil {
					log.Printf("error == %s",err )
				
			}
		}
	}

	return nil
}
