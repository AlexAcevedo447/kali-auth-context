#!/bin/sh
# Genera el par de claves RSA (PKCS#1) que el emisor de JWT (RS256) necesita
# para desarrollo local. Ver internal/infrastructure/security/jwt_access_token_issuer.go,
# que exige exactamente este formato (x509.ParsePKCS1PrivateKey / ParsePKCS1PublicKey).
# No regenera las claves si ya existen.
set -eu

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SECRETS_DIR="$REPO_ROOT/.secrets"
PRIVATE_KEY="$SECRETS_DIR/jwt_private.pem"
PUBLIC_KEY="$SECRETS_DIR/jwt_public.pem"

mkdir -p "$SECRETS_DIR"

if [ -f "$PRIVATE_KEY" ] && [ -f "$PUBLIC_KEY" ]; then
  echo "Las claves JWT de desarrollo ya existen en $SECRETS_DIR, no se regeneran."
  exit 0
fi

openssl genrsa -out "$PRIVATE_KEY" 2048 2>/dev/null
# `genrsa` puede emitir PKCS#8 ("BEGIN PRIVATE KEY") según la build de OpenSSL;
# el issuer usa x509.ParsePKCS1PrivateKey, que exige el formato tradicional.
openssl rsa -in "$PRIVATE_KEY" -traditional -out "$PRIVATE_KEY" 2>/dev/null
chmod 600 "$PRIVATE_KEY"

openssl rsa -in "$PRIVATE_KEY" -RSAPublicKey_out -out "$PUBLIC_KEY" 2>/dev/null
chmod 644 "$PUBLIC_KEY"

echo "Claves JWT de desarrollo generadas en $SECRETS_DIR"
