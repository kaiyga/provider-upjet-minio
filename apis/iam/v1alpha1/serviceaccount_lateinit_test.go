package v1alpha1

import "testing"

// Регрессия: late-init не должен писать expiration из tf state в spec.
// Бессрочный ключ MinIO отдаёт 1970-01-01, а refresh с таким expiration падает
// ("Expiration must between 15 minutes and 365 days"), из-за чего ресурс
// не наблюдается и не удаляется (зависший finalizer).
func TestServiceAccountLateInitializeSkipsExpiration(t *testing.T) {
	sa := &ServiceAccount{}
	changed, err := sa.LateInitialize([]byte(`{"expiration":"1970-01-01T00:00:00Z","target_user":"u"}`))
	if err != nil {
		t.Fatal(err)
	}
	if sa.Spec.ForProvider.Expiration != nil {
		t.Fatalf("expiration was late-initialized: %q", *sa.Spec.ForProvider.Expiration)
	}
	if !changed || sa.Spec.ForProvider.TargetUser == nil {
		t.Fatalf("other fields must still be late-initialized (changed=%v)", changed)
	}
}
