# Maintainer: ybenel <root@ybenel.cf>
pkgname=Neptune
pkgver=1.0.1
pkgrel=1
pkgdesc="Neptune a superfast mechanical keyboard sound app"
arch=('x86_64')
url="https://github.com/M1ndo/Neptune"
license=('AGPLv3')
depends=('make')

source=("https://github.com/M1ndo/Neptune/releases/download/v$pkgver/Neptune.tar.xz")
sha256sums=('44b045c049786265838d096d524c0b3e616f4d4f2f954e41d5277276660e3a63')

build() {
  cd "$srcdir"
  tar xf "$srcdir/$pkgname.tar.xz"
}

package() {
  cd "$srcdir"
  make DESTDIR="$pkgdir/" install
}