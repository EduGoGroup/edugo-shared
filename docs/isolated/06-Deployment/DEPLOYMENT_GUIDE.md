# Deployment - spec-04
## Publicar Release
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared

# Tag release desde dev
git checkout dev
git tag v0.7.0
git push origin v0.7.0

# Proyectos importan
go get github.com/EduGoGroup/edugo-shared/logger@v0.7.0
```
