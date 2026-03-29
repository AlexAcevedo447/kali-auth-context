package ports

type IUUIDProvider interface {
	Generate() string
}
