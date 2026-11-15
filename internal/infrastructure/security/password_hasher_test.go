package security

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Las contraseñas nuevas se guardan con Argon2id, que usa la contraseña
// completa: no hay truncamiento a 72 bytes como en bcrypt.
func TestArgon2id_UsesWholePasswordBeyond72Bytes(t *testing.T) {
	h := NewArgon2idHasher()
	long := "Ñandú#Clave1" + strings.Repeat("ñ", 50) // 114 bytes
	hash, err := h.Hash(long)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf("hash nuevo debe ser Argon2id, fue %q", hash[:10])
	}
	if err := h.Compare(hash, long); err != nil {
		t.Fatalf("la contraseña completa debe verificar: %v", err)
	}
	if err := h.Compare(hash, long[:72]+"x"); err == nil {
		t.Fatal("una contraseña que solo coincide en los primeros 72 bytes no debe verificar")
	}
	if h.NeedsRehash(hash) {
		t.Fatal("un hash Argon2id vigente no necesita rehash")
	}
}

// bcrypt solo se usa para verificar hashes heredados ("$2…"): su límite
// efectivo es de 72 bytes (lo que sigue no se compara). El login señala esos
// hashes con needs_rehash; la migración automática no está implementada.
func TestLegacyBcrypt_Effective72ByteLimitAndRehash(t *testing.T) {
	p72 := strings.Repeat("a", 72)
	if _, err := bcrypt.GenerateFromPassword([]byte(p72+"b"), bcrypt.MinCost); err == nil {
		t.Fatal("bcrypt debe rechazar crear hashes de más de 72 bytes")
	}
	legacy, err := bcrypt.GenerateFromPassword([]byte(p72), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	h := NewArgon2idHasher()
	if err := h.Compare(string(legacy), p72+"OTRO"); err != nil {
		t.Fatalf("comportamiento documentado de bcrypt: ignora lo que pasa de 72 bytes (%v)", err)
	}
	if !h.NeedsRehash(string(legacy)) {
		t.Fatal("un hash bcrypt heredado debe señalarse con needs_rehash")
	}
}
