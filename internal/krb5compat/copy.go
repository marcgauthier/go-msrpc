// Package krb5compat translates legacy Kerberos credential structures into
// the supported oiweiwei/gokrb5.fork/v9 structures without importing gokrb5/v8.
//
// This adapter copies exported matching fields, recursively translating nested
// structs, pointers, arrays, and slices. It does not use unsafe operations.
package krb5compat

import (
    "fmt"
    "reflect"
)

// Copy converts the contents of src into a non-nil pointer dst. Src may have a
// different named struct type (e.g., a legacy Kerberos v8 credential).
func Copy(dst, src any) error {
    d := reflect.ValueOf(dst)
    if !d.IsValid() || d.Kind() != reflect.Pointer || d.IsNil() {
        return fmt.Errorf("destination must be a non-nil pointer")
    }
    s := reflect.ValueOf(src)
    if !s.IsValid() || ((s.Kind() == reflect.Pointer || s.Kind() == reflect.Interface) && s.IsNil()) {
        return fmt.Errorf("source is nil")
    }
    return copyValue(d.Elem(), s)
}

func copyValue(dst, src reflect.Value) error {
    if !src.IsValid() {
        return fmt.Errorf("missing source value for %s", dst.Type())
    }
    if src.Kind() == reflect.Interface {
        if src.IsNil() { return nil }
        return copyValue(dst, src.Elem())
    }
    if dst.Kind() == reflect.Interface {
        if src.Type().AssignableTo(dst.Type()) {
            dst.Set(src)
            return nil
        }
        return fmt.Errorf("cannot assign %s to %s", src.Type(), dst.Type())
    }
    if src.Type().AssignableTo(dst.Type()) {
        dst.Set(src)
        return nil
    }
    if src.Kind() == reflect.Pointer {
        if src.IsNil() { return nil }
        if dst.Kind() == reflect.Pointer {
            if dst.IsNil() { dst.Set(reflect.New(dst.Type().Elem())) }
            return copyValue(dst.Elem(), src.Elem())
        }
        return copyValue(dst, src.Elem())
    }
    if dst.Kind() == reflect.Pointer {
        if dst.IsNil() { dst.Set(reflect.New(dst.Type().Elem())) }
        return copyValue(dst.Elem(), src)
    }
    switch dst.Kind() {
    case reflect.Struct:
        if src.Kind() != reflect.Struct {
            return fmt.Errorf("expected struct for %s, got %s", dst.Type(), src.Type())
        }
        for i := 0; i < dst.NumField(); i++ {
            f := dst.Type().Field(i)
            to := dst.Field(i)
            if !f.IsExported() || !to.CanSet() { continue }
            from := src.FieldByName(f.Name)
            if !from.IsValid() { continue } // v9-only optional field
            if err := copyValue(to, from); err != nil {
                return fmt.Errorf("%s.%s: %w", dst.Type(), f.Name, err)
            }
        }
        return nil
    case reflect.Slice:
        if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
            return fmt.Errorf("expected slice for %s, got %s", dst.Type(), src.Type())
        }
        if src.Kind() == reflect.Slice && src.IsNil() { return nil }
        result := reflect.MakeSlice(dst.Type(), src.Len(), src.Len())
        for i := 0; i < src.Len(); i++ {
            if err := copyValue(result.Index(i), src.Index(i)); err != nil {
                return fmt.Errorf("[%d]: %w", i, err)
            }
        }
        dst.Set(result)
        return nil
    case reflect.Array:
        if src.Kind() != reflect.Array && src.Kind() != reflect.Slice {
            return fmt.Errorf("expected array for %s, got %s", dst.Type(), src.Type())
        }
        if dst.Len() != src.Len() {
            return fmt.Errorf("array length mismatch: %d vs %d", dst.Len(), src.Len())
        }
        for i := 0; i < dst.Len(); i++ {
            if err := copyValue(dst.Index(i), src.Index(i)); err != nil {
                return fmt.Errorf("[%d]: %w", i, err)
            }
        }
        return nil
    case reflect.Map:
        if src.Kind() != reflect.Map { break }
        if src.Type().ConvertibleTo(dst.Type()) {
            dst.Set(src.Convert(dst.Type()))
            return nil
        }
    }
    if src.Type().ConvertibleTo(dst.Type()) {
        dst.Set(src.Convert(dst.Type()))
        return nil
    }
    return fmt.Errorf("cannot convert %s to %s", src.Type(), dst.Type())
}
