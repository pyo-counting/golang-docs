# helm-ecr-api

ECR에 저장된 Helm Chart를 조회하는 REST API 서버입니다.

## 실행 방법

1.  **AWS 자격 증명 설정**:
    ECR에 접근할 수 있는 AWS 자격 증명을 환경 변수로 설정합니다.

    ```sh
    export AWS_REGION="ap-northeast-2"
    export AWS_ACCESS_KEY_ID="YOUR_AWS_ACCESS_KEY"
    export AWS_SECRET_ACCESS_KEY="YOUR_AWS_SECRET_KEY"
    ```

2.  **의존성 설치**:

    ```sh
    go mod tidy
    ```

3.  **서버 실행**:

    ```sh
    go run ./cmd/api
    ```

    또는 Docker를 사용하여 실행할 수도 있습니다.

    ```sh
    docker run -p 8080:8080 \
      -e AWS_REGION="ap-northeast-2" \
      -e AWS_ACCESS_KEY_ID="YOUR_AWS_ACCESS_KEY" \
      -e AWS_SECRET_ACCESS_KEY="YOUR_AWS_SECRET_KEY" \
      helm-ecr-api:latest
    ```

## 설정

-   `PORT`: 서버가 실행될 포트를 지정합니다. (기본값: `8080`)
-   `HELM_REPOSITORIES`: API를 통해 노출할 ECR 리포지토리 목록을 콤마(`,`)로 구분하여 지정합니다. (필수)
    -   예: `export HELM_REPOSITORIES="my-charts/app1,my-charts/app2"`

## API 테스트

`curl`을 사용하여 API를 테스트할 수 있습니다. `repo`와 `tag` 파라미터를 실제 ECR에 있는 차트 정보로 변경하세요.

- **모든 Helm 차트 리포지토리 목록 조회**:
  ```sh
  curl "http://localhost:8080/v1/helm-charts"
  ```

- **리포지토리의 모든 차트 버전(태그) 조회**:
  ```sh
  curl "http://localhost:8080/v1/helm-charts/my-helm-charts/my-app"
  ```

- **특정 태그의 차트 정보 조회**:
  ```sh
  curl "http://localhost:8080/v1/helm-charts/my-helm-charts/my-app?tag=1.2.3"
  ```

- **특정 다이제스트의 차트 정보 조회**:
  ```sh
  curl "http://localhost:8080/v1/helm-charts/my-helm-charts/my-app?digest=sha256:..."
  ```

- **`values.yaml` 파일 내용 조회**:
  ```sh
  curl "http://localhost:8080/v1/helm-charts/my-helm-charts/my-app/files/values.yaml?tag=1.2.3"
  ```

- **`Chart.yaml` 파일 내용 조회**:
  ```sh
  curl "http://localhost:8080/v1/helm-charts/my-helm-charts/my-app/files/Chart.yaml?tag=1.2.3"
  ```

- **`values.schema.json` 파일 내용 조회**:
  ```sh
  curl "http://localhost:8080/v1/helm-charts/my-helm-charts/my-app/files/values.schema.json?tag=1.2.3"
  ```
