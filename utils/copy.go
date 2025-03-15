package utils

import "github.com/jinzhu/copier"

func DeepClone(src any, dest any) error {
	return copier.Copy(dest, src)
}
