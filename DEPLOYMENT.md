# Despliegue a AWS

Este proyecto está repartido en 5 repositorios independientes bajo
`github.com/ntorres1986`:

| Repo | Contenido | Destino |
|---|---|---|
| `nexed-shared` | Modelos de datos compartidos (módulo Go, público) | — (solo dependencia de los 2 backends) |
| `nexed-customer-backend` | API de la app de colegios/padres | Elastic Beanstalk |
| `nexed-customer-frontend` | Frontend de colegios/padres | S3 + CloudFront |
| `nexed-admin-backend` | API del panel de super-admin | Elastic Beanstalk |
| `nexed-admin-frontend` | Frontend del panel de super-admin | S3 + CloudFront |

Cada uno de los 4 repos de apps tiene su propio `.github/workflows/deploy.yml`
que se dispara con **push a `dev`** y despliega directo al único ambiente:
**`production`** (no hay staging — es un setup de un solo ambiente, el
nombre "production" es solo el nombre del GitHub Environment y de los
recursos de AWS, no tiene que ver con qué rama lo dispara).

## Estado actual de la infraestructura (ya creada)

- **RDS MySQL** `nexed-db` — `db.t3.micro`, base `matriculas`, contraseña
  gestionada por AWS Secrets Manager (managed master user password).
- **2 apps de Elastic Beanstalk** (plataforma Go, tier "Single instance"):
  `nexed-customer` (env `nexed-customer-prod`) y `nexed-admin` (env
  `nexed-admin-prod`).
- **2 buckets S3**: `nexed-customer-frontend-prod`, `nexed-admin-frontend-prod`
  (acceso público bloqueado — se sirven vía CloudFront con Origin Access
  Control).
- **Rol IAM** `nexed-github-deploy` — asumido por GitHub Actions vía OIDC,
  confianza restringida al Environment `production` de cada uno de los 4
  repos (usando los IDs inmutables de GitHub en el `sub` del token, no solo
  el nombre — ver nota más abajo).
- **CloudFront**: pendiente — la cuenta de AWS necesitó verificación manual
  vía un caso de soporte antes de poder crear distribuciones. Una vez
  resuelto, faltan las 2 distribuciones (una por bucket) con su Origin
  Access Control.

### Nota sobre el `sub` del token OIDC

GitHub incluye IDs numéricos inmutables en el claim `sub`, no solo los
nombres: `repo:ntorres1986@44281748/nexed-customer-backend@1366956460:environment:production`
en vez del clásico `repo:ntorres1986/nexed-customer-backend:environment:production`.
Si algún día hay que tocar la trust policy del rol, hay que usar esos IDs
(se consultan con `gh api repos/ntorres1986/<repo> --jq .id` y `gh api user
--jq .id`), no solo el nombre — si no, `AssumeRoleWithWebIdentity` falla con
"Not authorized" aunque el nombre del repo esté bien escrito.

## Variables de entorno en cada uno de los 2 EB (`nexed-customer-prod`, `nexed-admin-prod`)

Ya cargadas: `APP_PORT`, `DB_HOST`, `DB_NAME`, `DB_PORT`, `DB_USER`,
`DB_PASSWORD`, `JWT_SECRET`, `UPLOADS_DIR`.

Pendiente: `FRONTEND_URL` (URL pública del `nexed-customer-frontend`, se
carga una vez exista su CloudFront) y, en `nexed-customer-backend`, las
`PLACETOPAY_*` cuando haya credenciales reales de producción de AvalPay.

⚠️ **Revisa el límite de tamaño de subida de nginx en EB** — la plataforma
Go trae un `client_max_body_size` por defecto bajo (1MB), y este proyecto
sube logos y plantillas PDF. Si ves errores 413 al subir archivos, hay que
agregar un `.platform/nginx/conf.d/uploads.conf` en el repo con
`client_max_body_size 10M;`.

⚠️ **Pendiente, fuera de alcance de esta configuración de CI/CD**: hoy
`nexed-admin-backend` y `nexed-customer-backend` comparten uploads mediante
una carpeta local en disco en desarrollo. En AWS cada backend corre en su
propia instancia EB — máquinas separadas — así que esa ruta relativa no
sirve en producción. Para que la subida de logos institucionales funcione
en la nube hace falta mover esos uploads a S3 (cambio de código, no
incluido aquí).

## GitHub: secrets y variables (Environment `production`, en cada uno de los 4 repos)

**Secret**: `AWS_DEPLOY_ROLE_ARN` = `arn:aws:iam::341853291054:role/nexed-github-deploy`

**Variables**, ya cargadas:
- `AWS_REGION` = `us-east-1`
- `CUSTOMER_BACKEND_EB_APP` / `CUSTOMER_BACKEND_EB_ENV` (repo customer-backend)
- `ADMIN_BACKEND_EB_APP` / `ADMIN_BACKEND_EB_ENV` (repo admin-backend)
- `CUSTOMER_FRONTEND_BUCKET`, `CUSTOMER_API_URL` (repo customer-frontend)
- `ADMIN_FRONTEND_BUCKET`, `ADMIN_API_URL` (repo admin-frontend)

Pendiente cuando exista CloudFront: `CUSTOMER_FRONTEND_CF_DISTRIBUTION_ID`,
`ADMIN_FRONTEND_CF_DISTRIBUTION_ID`.

## Primer despliegue / relanzar uno

Cada workflow corre solo en push a `dev`. Para forzar uno sin tocar código:
`gh workflow run deploy.yml --repo ntorres1986/<repo> --ref dev` (o desde la
pestaña Actions de GitHub).

## nexed-shared (módulo Go compartido)

`nexed-customer-backend` y `nexed-admin-backend` dependen de
`github.com/ntorres1986/nexed-shared` como un módulo Go normal. Si cambias
algo ahí:

1. Commit + push a `nexed-shared` (rama `dev`).
2. Tag de una versión **nueva** (ej. `git tag v0.1.2 && git push origin
   v0.1.2`) — nunca reuses un tag existente: `proxy.golang.org` cachea
   resultados por versión y, si el repo era privado cuando alguien pidió
   esa versión por primera vez, el caché negativo queda pegado a ese tag
   para siempre (nos pasó con `v0.1.0` — tuvimos que saltar a `v0.1.1`).
3. En cada backend: `go get github.com/ntorres1986/nexed-shared@vX.Y.Z &&
   go mod tidy`, commit del `go.mod`/`go.sum`, push.
