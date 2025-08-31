<h2 align="center">Back database</h3>
  <p align="center">
    <a href="https://github.com/TheoMeunier/Filesox/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/TheoMeunier/Filesox/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>

## About The Project

The tool allows you to dump your Postgres database and then store it in an S3 bucket for as long as you want.

## Getting Started

#### Docker

1. Create a docker-compose file

```yml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: backup-database
  namespace: production
  labels:
    app: backup
    component: database
spec:
  schedule: "30 2 * * *"
  timeZone: "Europe/Paris"
  concurrencyPolicy: Forbid
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 1
  startingDeadlineSeconds: 300
  suspend: false

  jobTemplate:
    spec:
      template:
        metadata:
          labels:
            app: backup
            component: database
        spec:
          restartPolicy: OnFailure

          containers:
            - name: backup-container-postgres
              image: ghcr.io/theomeunier/worker-backup-db/postgres:latest
              envFrom:
                - secretRef:
                    name: backup-database-secrets
              resources:
                requests:
                  memory: "128Mi"
                  cpu: "100m"
                limits:
                  memory: "512Mi"
                  cpu: "500m"
```

### 2. Configure the `.env` file

#### 2.1 DATABASE Configuration:

- `DATABASE_USER` : The username of your database
- `DATABASE_PASSWORD` : The password of your database
- `DATABASE_HOST`: The URL of your database
- `DATABASE_NAME`: The name of your database
- `DATABASE_PORT`: The port of your database

#### 2.2 S3 Configuration:

- `S3_ENDPOINT_URL` : Endpoint of your S3 provider
- `S3_REGION` : Region of your S3 provider
- `S3_REGION` : The region of your S3 provider
- `S3_ACCESS_KEY` : Access key of your S3 provider
- `S3_SECRET_KEY` : Secret key of your S3 provider

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any
contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also
simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.