# Copilot Instructions for DevOps/SRE Engineer

## 역할 정보 (Role Information)
- **직무**: DevOps/SRE 엔지니어
- **Go 버전**: Go 1.24.3
- **경험 수준**: Go 입문자

## 코딩 가이드라인 (Coding Guidelines)

### 1. Go 공식 문서 기반 Best Practices
- 항상 [Go 공식 문서](https://golang.org/doc/)의 best practices를 따라주세요
- [Effective Go](https://golang.org/doc/effective_go.html) 가이드라인을 준수해주세요
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)를 참조해주세요

### 2. 코드 제안 시 요구사항
- **입문자 친화적**: 복잡한 개념을 단계별로 설명
- **명확한 주석**: 코드의 각 부분이 무엇을 하는지 설명
- **에러 처리**: Go의 관용적인 에러 처리 패턴 사용
- **테스트 코드**: 가능한 경우 테스트 예제도 함께 제공

### 3. DevOps/SRE 관련 우선순위
- **로깅**: structured logging 사용 (예: `log/slog` 패키지)
- **메트릭**: Prometheus 메트릭 패턴
- **설정 관리**: 환경변수, 설정 파일 best practices
- **HTTP 서버**: graceful shutdown, health checks
- **Docker**: 효율적인 Dockerfile 작성법
- **Kubernetes**: Go 애플리케이션의 k8s 배포 패턴

### 4. 참조 문서 포함 요구사항
코드 제안 시 다음 참조 문서를 함께 제공해주세요:
- 관련 Go 공식 문서 링크
- 패키지 문서 (pkg.go.dev 링크)
- 사용된 패턴이나 개념의 설명 링크
- DevOps/SRE 관련 best practice 문서

### 5. 코드 스타일
- `gofmt`로 포맷된 코드
- 전체적인 코드 스타일 가이드: [uber-go/guide](https://github.com/uber-go/guide/blob/master/style.md)
- 명명 규칙: Go 공식 가이드라인 준수
- 패키지 구조: 표준 Go 프로젝트 레이아웃
    - 참고 구조: [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- 의존성 관리: Go modules 사용

### 6. 예제 템플릿
코드 제안 시 다음과 같은 구조로 제공해주세요:

```go
// 간단한 설명
// 참조: [관련 문서 링크]

package main

import (
    // 필요한 패키지들
)

// 주요 함수나 구조체에 대한 설명
func main() {
    // 단계별 주석과 함께 코드 작성
}
```

**참조 문서:**
- [Go 공식 문서](https://golang.org/doc/)
- [패키지 문서](https://pkg.go.dev/)
- [Go best practices 가이드](https://github.com/golang/go/wiki/CodeReviewComments)

### 7. 보안 및 컴플라이언스
- **시크릿 관리**: 환경변수, Vault, Kubernetes secrets 패턴
- **TLS/암호화**: 안전한 통신 구현
- **인증/인가**: JWT, OAuth2, RBAC 패턴
- **입력 검증**: 보안 취약점 방지를 위한 입력 유효성 검사
- **참조**: [Go Security Checklist](https://github.com/securecodewarrior/go-security-checklist)

### 8. 모니터링 및 관찰가능성 (Observability)
- **메트릭**: Prometheus, Grafana 통합
- **로깅**: 구조화된 로깅, 로그 레벨 관리
- **트레이싱**: OpenTelemetry, Jaeger 통합
- **헬스체크**: `/health`, `/ready`, `/metrics` 엔드포인트
- **참조**: [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)

### 9. 성능 및 확장성
- **동시성**: 고루틴, 채널 패턴 best practices
- **메모리 관리**: 메모리 누수 방지, 프로파일링
- **데이터베이스**: 커넥션 풀, 쿼리 최적화
- **캐싱**: Redis, in-memory 캐시 패턴
- **참조**:
    - [Go Performance](https://github.com/dgryski/go-perfbook)
    - [Go Optimization Guide](https://goperf.dev/)
    - [Scaling GO Applications](https://betterstack.com/community/guides/scaling-go/)

### 10. 인프라스트럭처 as Code
- **Terraform**: Go로 Terraform provider 작성
- **Kubernetes Operators**: controller-runtime 사용
- **CI/CD**: GitHub Actions, GitLab CI 파이프라인
- **컨테이너**: 멀티스테이지 빌드, 보안 이미지 작성
- **참조**: [Kubernetes Operator SDK](https://sdk.operatorframework.io/)

### 11. 운영 환경 고려사항
- **설정 관리**: 12-factor app 원칙 준수
- **로그 로테이션**: 디스크 공간 관리
- **리소스 제한**: CPU, 메모리 제한 설정
- **백업/복구**: 데이터 백업 전략
- **장애 대응**: Circuit breaker, retry 패턴

### 12. 팀 협업 및 문서화
- **코드 리뷰**: 효과적인 PR 작성법
- **문서화**: README, API 문서, 운영 가이드
- **버전 관리**: semantic versioning, changelog 관리
- **의존성 관리**: go.mod 보안 업데이트

---

이 지침을 따라 Go 입문자인 DevOps/SRE 엔지니어에게 최적화된 코드와 설명을 제공해주세요.