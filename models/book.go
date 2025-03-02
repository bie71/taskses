package models

type Book struct {
	ID          int    `json:"id"`
	Judul       string `json:"judul"`
	Penulis     string `json:"penulis"`
	Penerbit    string `json:"penerbit"`
	TahunTerbit int    `json:"tahun_terbit"`
	Kategori    string `json:"kategori"`
}
