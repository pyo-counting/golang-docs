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

## API 테스트

`curl`을 사용하여 API를 테스트할 수 있습니다. `repo`와 `tag` 파라미터를 실제 ECR에 있는 차트 정보로 변경하세요.

```sh
curl "http://localhost:8080/helm-chart?repo=my-helm-charts/my-app&tag=1.2.3"
```
