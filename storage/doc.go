// Package storage define interfaces para operaciones de almacenamiento de archivos.
//
// Proporciona dos interfaces principales:
//   - [Client] para operaciones CRUD (Download, Upload, Delete, Exists, GetMetadata)
//   - [PresignClient] para generacion de URLs pre-firmadas (Upload/Download)
//
// La implementacion S3 esta en el sub-modulo [storage/s3].
package storage
