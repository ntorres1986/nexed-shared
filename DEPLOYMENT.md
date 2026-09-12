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
que se dispara solo con push a `dev` o `main`:

- **push a `dev`** → despliega al ambiente **staging**
- **push a `main`** → despliega al ambiente **production**

Ninguno funciona todavía: primero hay que crear la infraestructura en AWS y
configurar los **GitHub Environments** "staging"/"production" de cada repo
(con sus secrets/variables). Esto no lo puedo hacer yo desde aquí (no tengo
credenciales de AWS ni acceso a la consola de GitHub) — son los pasos
manuales que faltan.

## 1. Recursos a crear en AWS (una sola vez, duplicados por ambiente salvo la base de datos)

- **4 buckets S3** — uno por frontend y por ambiente:
  `nexed-customer-frontend-staging`, `nexed-customer-frontend-prod`,
  `nexed-admin-frontend-staging`, `nexed-admin-frontend-prod` (los nombres
  son libres). Sin "static website hosting" activado: CloudFront los sirve
  como origen privado con OAC.
- **4 distribuciones CloudFront** — una por bucket, plan gratuito de $0/mes.
  Origin con Origin Access Control (OAC), bucket no público.
- **1 instancia RDS MySQL** — compartida por *ambos ambientes y ambos
  backends*: staging y producción usan la MISMA base de datos `matriculas`
  (decisión explícita del proyecto: no hay una base separada por ambiente).
  `db.t3.micro`, single-AZ.
- **4 entornos de Elastic Beanstalk**, plataforma **Go**, tier
  **"Single instance"** (no "Load balanced"):
  - `nexed-customer-backend` → app `nexed-customer` / envs
    `nexed-customer-staging` y `nexed-customer-prod`
  - `nexed-admin-backend` → app `nexed-admin` / envs `nexed-admin-staging`
    y `nexed-admin-prod`

### Variables de entorno a configurar EN CADA ENTORNO DE EB

Esto se configura directo en la consola de EB (Configuration → Software →
Environment properties), **no** en GitHub — son configuración del
servidor, no del pipeline. **Todos los entornos (staging y prod) apuntan a
la misma RDS/base de datos.**

Para ambos backends, en los 4 entornos:

```
APP_ENV=production      # o "staging" si quieres distinguirlo en logs/emails
APP_PORT=5000            # EB (plataforma Go) enruta nginx hacia el puerto 5000 por defecto
DB_HOST=<endpoint de tu RDS>
DB_PORT=3306
DB_NAME=matriculas
DB_USER=...
DB_PASSWORD=...
JWT_SECRET=<un secreto largo, igual en ambos backends si comparten sesión>
FRONTEND_URL=https://<dominio de CloudFront del nexed-customer-frontend de ESTE ambiente>
UPLOADS_DIR=/var/app/current/uploads
```

`nexed-customer-backend` además necesita las variables de `PLACETOPAY_*` (ver
`.env.example` del repo) y, si quieres el webhook de AvalPay funcionando,
`PLACETOPAY_NOTIFICATION_URL` apuntando a la URL pública del propio backend
de ese ambiente.

⚠️ **Revisa el límite de tamaño de subida de nginx en EB** — la plataforma
Go trae un `client_max_body_size` por defecto bajo (1MB), y este proyecto
sube logos y plantillas PDF. Si ves errores 413 al subir archivos, hay que
agregar un `.platform/nginx/conf.d/uploads.conf` en el repo con
`client_max_body_size 10M;`.

⚠️ **Pendiente, fuera de alcance de esta configuración de CI/CD**: hoy
`nexed-admin-backend` y `nexed-customer-backend` comparten uploads mediante
una carpeta local en disco (`UPLOADS_DIR=../nexed-customer-backend/uploads`
en desarrollo). En AWS cada backend corre en su propia instancia EB —
máquinas separadas — así que esa ruta relativa no sirve en producción. Para
que la subida de logos institucionales funcione en la nube hace falta mover
esos uploads a S3 (cambio de código en ambos backends, no incluido aquí).

## 2. Rol de IAM para GitHub Actions (OIDC — sin llaves de AWS guardadas en GitHub)

1. En IAM → Identity providers, agrega un proveedor OIDC:
   - URL: `https://token.actions.githubusercontent.com`
   - Audience: `sts.amazonaws.com`
2. Crea un rol con esta *trust policy* (reemplaza `<ACCOUNT_ID>` y
   `<REPO>` por cada uno de los 4 repos de apps — puedes usar un solo rol
   compartido si el `sub` permite los 4, o un rol por repo):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::<ACCOUNT_ID>:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": [
            "repo:ntorres1986/<REPO>:ref:refs/heads/dev",
            "repo:ntorres1986/<REPO>:ref:refs/heads/main"
          ]
        }
      }
    }
  ]
}
```

Si prefieres condicionar por GitHub Environment en vez de por rama (más
preciso, recomendado por GitHub cuando usás Environments), usa en su lugar:

```json
"token.actions.githubusercontent.com:sub": [
  "repo:ntorres1986/<REPO>:environment:staging",
  "repo:ntorres1986/<REPO>:environment:production"
]
```

3. Adjunta una política de permisos. La más simple es usar
   `AWSElasticBeanstalkFullAccess` (manejada por AWS) más esto para S3/CloudFront:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "FrontendBuckets",
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:GetObject", "s3:ListBucket", "s3:DeleteObject"],
      "Resource": [
        "arn:aws:s3:::<bucket-nexed-customer-frontend-staging>",
        "arn:aws:s3:::<bucket-nexed-customer-frontend-staging>/*",
        "arn:aws:s3:::<bucket-nexed-customer-frontend-prod>",
        "arn:aws:s3:::<bucket-nexed-customer-frontend-prod>/*",
        "arn:aws:s3:::<bucket-nexed-admin-frontend-staging>",
        "arn:aws:s3:::<bucket-nexed-admin-frontend-staging>/*",
        "arn:aws:s3:::<bucket-nexed-admin-frontend-prod>",
        "arn:aws:s3:::<bucket-nexed-admin-frontend-prod>/*"
      ]
    },
    {
      "Sid": "CloudFrontInvalidation",
      "Effect": "Allow",
      "Action": ["cloudfront:CreateInvalidation"],
      "Resource": "*"
    }
  ]
}
```

Copia el ARN del rol resultante (`arn:aws:iam::<ACCOUNT_ID>:role/...`) — lo
necesitas en el siguiente paso.

## 3. GitHub Environments, secrets y variables (por cada uno de los 4 repos de apps)

En cada repo → Settings → **Environments**, crea dos: `staging` y
`production`. Dentro de cada Environment vas a cargar los mismos *nombres*
de secret/variable, pero con el valor que corresponde a ese ambiente (por
ejemplo, `CUSTOMER_BACKEND_EB_ENV=nexed-customer-staging` en el Environment
"staging" y `nexed-customer-prod` en "production").

**Secrets** (sensible), en cada Environment:
- `AWS_DEPLOY_ROLE_ARN` — el ARN del rol del paso 2.

**Variables** (no sensible), en cada Environment — solo las que aplican a
ese repo:
- `AWS_REGION` — ej. `us-east-1` (puede ser la misma en ambos ambientes)
- `CUSTOMER_API_URL` / `ADMIN_API_URL` — URL pública del backend de ESE
  ambiente
- `CUSTOMER_FRONTEND_BUCKET` / `CUSTOMER_FRONTEND_CF_DISTRIBUTION_ID`
- `ADMIN_FRONTEND_BUCKET` / `ADMIN_FRONTEND_CF_DISTRIBUTION_ID`
- `CUSTOMER_BACKEND_EB_APP` / `CUSTOMER_BACKEND_EB_ENV`
- `ADMIN_BACKEND_EB_APP` / `ADMIN_BACKEND_EB_ENV`

## 4. Primer despliegue

Cada workflow corre automáticamente en cada push a `dev` o `main`. Para
forzar uno sin tocar código, cada uno tiene `workflow_dispatch` — se puede
lanzar a mano desde la pestaña **Actions** de cada repo (elige la rama antes
de lanzarlo, ya que determina el ambiente).

## nexed-shared (módulo Go compartido)

`nexed-customer-backend` y `nexed-admin-backend` dependen de
`github.com/ntorres1986/nexed-shared` como un módulo Go normal (no vive más
dentro de ninguno de los 2 backends). Si cambias algo en `nexed-shared`:

1. Commit + push a `nexed-shared` (rama `dev` para cambios en curso).
2. Tag de una nueva versión (ej. `git tag v0.1.2 && git push origin v0.1.2`)
   — **usa siempre un tag nuevo**, nunca reuses uno existente: el proxy de
   Go (`proxy.golang.org`) cachea resultados por versión, y si el tag se
   reutiliza sobre un commit distinto el caché puede quedar desincronizado.
3. En cada backend: `go get github.com/ntorres1986/nexed-shared@vX.Y.Z &&
   go mod tidy`, commit del `go.mod`/`go.sum` actualizado, push.
