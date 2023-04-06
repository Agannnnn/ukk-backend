package utils

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Petugas struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_petugas" form:"id_petugas"`
	NamaDepan    string         `gorm:"type:varchar(100)" json:"nama_depan" form:"nama_depan"`
	NamaBelakang string         `gorm:"type:varchar(100)" json:"nama_belakang" form:"nama_belakang"`
	Username     string         `gorm:"type:varchar(50)" json:"username" form:"username"`
	Password     string         `gorm:"type:varchar(50)" json:"password" form:"password"`
	Level        string         `gorm:"type:level" json:"level" form:"level"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	LelangDibuat  []Lelang `gorm:"foreignKey:PetugasPembuatID"`
	LelangDitutup []Lelang `gorm:"foreignKey:PetugasPenutupID"`
}

type Masyarakat struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_masyarakat" form:"id_masyarakat"`
	NamaDepan    string    `gorm:"type:varchar(100)" json:"nama_depan" form:"nama_depan"`
	NamaBelakang string    `gorm:"type:varchar(100)" json:"nama_belakang" form:"nama_belakang"`
	Username     string    `gorm:"type:varchar(50)" json:"username" form:"username"`
	Password     string    `gorm:"type:varchar(50)" json:"password" form:"password"`
	NoTelp       string    `gorm:"type:varchar(14)" json:"no_telp" form:"no_telp"`
	Email        string    `gorm:"type:varchar(100)" json:"email" form:"email"`

	Penawaran []Penawaran `gorm:"foreignKey:MasyarakatID"`
	Langganan []Langganan `gorm:"foreignKey:MasyarakatID"`
}

type Barang struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_barang" form:"id_barang"`
	Nama      string    `gorm:"type:varchar(200)" json:"nama" form:"nama"`
	Deskripsi string    `json:"deskripsi" form:"deskripsi"`
	Diunggah  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"diunggah" form:"diunggah"`
	Deleted   gorm.DeletedAt

	FotoBarang []FotoBarang     `gorm:"foreignKey:BarangID"`
	Kategori   []DetailKategori `gorm:"foreignKey:BarangID"`
	Lelang     []Lelang         `gorm:"foreignKey:BarangID"`
}

type Lelang struct {
	ID               uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_lelang" form:"id_lelang"`
	MulaiLelang      time.Time `gorm:"type:timestamp" json:"mulai_lelang" form:"mulai_lelang"`
	SelesaiLelang    time.Time `gorm:"type:timestamp" json:"selesai_lelang" form:"selesai_lelang"`
	Timeout          int       `gorm:"type:integer" json:"timeout" form:"timeout"`
	Status           string    `gorm:"type:status_lelang;default:dibuka" json:"status_lelang" form:"status_lelang"`
	HargaAwal        int       `gorm:"type:integer" json:"harga_awal" form:"harga_awal"`
	MinPenawaran     int       `gorm:"type:integer" json:"min_penawaran" form:"min_penawaran"`
	BarangID         uuid.UUID `gorm:"type:uuid;" json:"id_barang" form:"id_barang"`
	PetugasPembuatID uuid.UUID `gorm:"type:uuid;" json:"id_petugas_pembuat" form:"id_petugas_pembuat"`
	PetugasPenutupID uuid.UUID `gorm:"type:uuid;default:null;" json:"id_petugas_penutup" form:"id_petugas_penutup"`

	Barang         Barang
	PetugasPembuat Petugas
	PetugasPenutup Petugas

	Penawaran []Penawaran `gorm:"foreignKey:LelangID"`
	Langganan []Langganan `gorm:"foreignKey:LelangID"`
}

type Penawaran struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_penawaran" form:"id_penawaran"`
	Harga        int       `gorm:"type:int" json:"harga" form:"harga"`
	Timestamp    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"timestamp" form:"timestamp"`
	LelangID     uuid.UUID `gorm:"type:uuid;" json:"id_lelang" form:"id_lelang"`
	MasyarakatID uuid.UUID `gorm:"type:uuid;" json:"id_masyarakat" form:"id_masyarakat"`

	Lelang     Lelang
	Masyarakat Masyarakat
}

type FotoBarang struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_foto_barang" form:"id_foto_barang"`
	Filename  string    `gorm:"type:varchar(70)" json:"filename" form:"filename"`
	Deskripsi string    `gorm:"type:text" json:"deskripsi" form:"deskripsi"`
	BarangID  uuid.UUID `gorm:"type:uuid;" json:"id_barang" form:"id_barang"`

	Barang Barang
}

type Kategori struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_kategori" form:"id_kategori"`
	Nama string    `gorm:"type:varchar(50)" json:"nama" form:"nama"`
}

type DetailKategori struct {
	BarangID   uuid.UUID `gorm:"type:uuid;" json:"id_barang" form:"id_barang"`
	KategoriID uuid.UUID `gorm:"type:uuid;" json:"id_kategori" form:"id_kategori"`

	Barang   Barang
	Kategori Kategori
}

type Langganan struct {
	LelangID     uuid.UUID `gorm:"type:uuid;" json:"id_lelang" form:"id_lelang"`
	MasyarakatID uuid.UUID `gorm:"type:uuid;" json:"id_masyarakat" form:"id_masyarakat"`

	Lelang     Lelang
	Masyarakat Masyarakat
}
