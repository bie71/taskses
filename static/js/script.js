        // Loop mencari ID yang cocok dalam daftarBuku.
        // Jika ditemukan, hapus elemen dari slice menggunakan append().
        // Jika tidak ditemukan, kirim 404 Not Found.
        //  Mengambil data buku dari backend menggunakan fetch API (GET /buku).
        //  Mengonversi response menjadi JSON (response.json()).
        //  Menghapus daftar lama sebelum memperbarui list (bukuList.innerHTML = "").
        //  Menggunakan loop untuk menampilkan daftar buku dalam format <li>:
        // Tampilkan informasi buku (judul - penulis (tahun_terbit)).
        // Tambahkan tombol "Edit" → Memanggil editBuku(id, judul, penulis, penerbit, tahunTerbit, kategori).
        // Tambahkan tombol "Hapus" → Memanggil hapusBuku(id).
        // Fungsi untuk memuat daftar buku
        async function loadBuku() {
            const token = localStorage.getItem("token");
            if (!token) {
                alert("Anda harus login terlebih dahulu");
                window.location.href = "/login"; // Redirect jika tidak ada token
            }
            const response = await fetch("/buku");
            const data = await response.json();
            let bukuList = document.getElementById("buku-list");
            bukuList.innerHTML = "";

            data.forEach(buku => {
                bukuList.innerHTML += `
            <li>
                    <div class="book-info">
                        <strong>${buku.judul}</strong> - ${buku.penulis} (${buku.tahun_terbit})
                    </div>
                    <div class="button-group">
                        <button class="edit" onclick="editBuku(${buku.id}, '${buku.judul}', '${buku.penulis}', '${buku.penerbit}', ${buku.tahun_terbit}, '${buku.kategori}')">Edit</button>
                        <button class="hapus" onclick="hapusBuku(${buku.id})">Hapus</button>
                    </div>
            </li>`;
            });
        }

        // Mengambil inputan dari form (document.getElementById().value).
        //  Menggunakan fetch untuk mengirim request POST /buku.
        //  Konversi data ke JSON (JSON.stringify()) sebelum dikirim.
        //  Reset form setelah data berhasil dikirim (form.reset()).
        //  Memanggil loadBuku() untuk memperbarui daftar buku.
        // Fungsi untuk menambah buku
        async function tambahBuku() {
            const judul = document.getElementById("judul").value;
            const penulis = document.getElementById("penulis").value;
            const penerbit = document.getElementById("penerbit").value;
            const tahunTerbit = document.getElementById("tahun_terbit").value;
            const kategori = document.getElementById("kategori").value;

            await fetch("/buku", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ 
                    judul, 
                    penulis, 
                    penerbit, 
                    tahun_terbit: parseInt(tahunTerbit), 
                    kategori 
                })
            });

            // Reset form dan reload daftar buku
            document.getElementById("form-buku").reset();
            loadBuku();
        }

        // ✅ Mengirim request DELETE ke backend (/buku/:id).
        // ✅ Memanggil loadBuku() agar daftar diperbarui setelah penghapusan.
        // Fungsi untuk menghapus buku
        async function hapusBuku(id) {
            await fetch(`/buku/${id}`, { method: "DELETE" });
            loadBuku();
        }

        // Fungsi untuk mengedit buku
        function editBuku(id, judul, penulis, penerbit, tahunTerbit, kategori) {
            document.getElementById("id").value = id;
            document.getElementById("judul").value = judul;
            document.getElementById("penulis").value = penulis;
            document.getElementById("penerbit").value = penerbit;
            document.getElementById("tahun_terbit").value = tahunTerbit;
            document.getElementById("kategori").value = kategori;
            document.getElementById("tambah-btn").style.display = "none";
            document.getElementById("update-btn").style.display = "inline";
        }
        
        // Mengisi form dengan data buku yang akan diedit.
        //  Menyembunyikan tombol "Tambah" dan menampilkan tombol "Perbarui"**.
        // Fungsi untuk memperbarui buku
        async function updateBuku() {
            const id = document.getElementById("id").value;
            const judul = document.getElementById("judul").value;
            const penulis = document.getElementById("penulis").value;
            const penerbit = document.getElementById("penerbit").value;
            const tahunTerbit = document.getElementById("tahun_terbit").value;
            const kategori = document.getElementById("kategori").value;

            await fetch(`/buku/${id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ 
                    id: parseInt(id), 
                    judul, 
                    penulis, 
                    penerbit, 
                    tahun_terbit: parseInt(tahunTerbit), 
                    kategori 
                })
            });

            // Reset form dan reload daftar buku
            document.getElementById("form-buku").reset();
            document.getElementById("tambah-btn").style.display = "inline";
            document.getElementById("update-btn").style.display = "none";
            loadBuku();
        }

        
        document.addEventListener("DOMContentLoaded", loadBuku);


function logout() {
            localStorage.removeItem("token"); // Hapus token dari localStorage
            window.location.href = "/login"; // Redirect ke halaman login
        }


async function updateUser() {
    const username = document.getElementById("new-username").value;
    const password = document.getElementById("new-password").value;
    const oldPassword = document.getElementById("old-password").value;


    const token = localStorage.getItem("token"); // Ambil token dari localStorage
     if (!token) {
        alert("Anda belum login!");
        window.location.href = "/login";
        return;
      }

      const updateButton = document.getElementById("update-btn");
        updateButton.innerText = "Memproses...";
        updateButton.disabled = true;

  try {
        const response = await fetch("/user", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + token
            },
            body: JSON.stringify({ 
                username, 
                new_password: password, 
                old_password: oldPassword 
            })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || "Terjadi kesalahan saat memperbarui user.");
        }

        alert("✅ " + data.message);
        window.location.reload(); // Refresh halaman setelah update berhasil
    } catch (error) {
        alert("❌ " + error.message);
    } finally {
        updateButton.innerText = "Update";
        updateButton.disabled = false;
    }
}
