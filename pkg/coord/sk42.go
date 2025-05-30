//nolint:all
package coord

import "math"

const (
	ro float64 = 180 / math.Pi * 3600 // Число угловых секунд в радиане
	// Эллипсоид Красовского.
	aP  float64 = 6378245         // Большая полуось
	alP float64 = 1 / 298.3       // Сжатие
	e2P float64 = 2*alP - alP*alP // Квадрат эксцентриситета

	// Эллипсоид WGS84 (GRS80, эти два эллипсоида сходны по большинству параметров).
	aW  float64 = 6378137           // Большая полуось
	alW float64 = 1 / 298.257223563 // Сжатие
	e2W float64 = 2*alW - alW*alW   // Квадрат эксцентриситета

	// Вспомогательные значения для преобразования эллипсоидов.
	a   float64 = (aP + aW) / 2
	e2  float64 = (e2P + e2W) / 2
	da  float64 = aW - aP
	de2 float64 = e2W - e2P

	// dx     dy      dz   wx wy    wz    ms
	// 23.92 -141.27 -80.9 0  0.35  0.82 -0.12*10-6 ГОСТ 51794-2001
	// 25    -141    -78.5 0 -0.35 -0.736 0         EPSG::15865
	dx float64 = 23.92
	dy float64 = -141.27
	dz float64 = -80.9

	// Угловые элементы трансформирования, в секундах.
	wx = 0
	wy = -0.35
	wz = -0.82
	// Дифференциальное различие масштабов.
	ms = -0.12e-6
)

func Wgs84_sk42(lat, lon, alt float64) (int, int, int) {
	return Sk42ll2Meters(wgs84_sk42ll(lat, lon, alt))
}

func Sk42_wgs(x, y int) (float64, float64) {
	lat, lon := sk42xy_to_sk42latlon(x, y)
	return sk42ll_wgs84(lat, lon, 0)
}

func wgs84_sk42ll(lat, lon, alt float64) (lat1, lon1 float64) {
	lat1 = lat - dB(lat, lon, alt)/3600
	lon1 = lon - dL(lat, lon, alt)/3600

	return
}

func sk42ll_wgs84(lat, lon, alt float64) (lat1, lon1 float64) {
	lat1 = lat + dB(lat, lon, alt)/3600
	lon1 = lon + dL(lat, lon, alt)/3600

	return
}

func dB(bd, ld, h float64) float64 {
	b := bd * math.Pi / 180
	l := ld * math.Pi / 180

	m := a * (1 - e2) / math.Pow(1-e2*math.Pow(math.Sin(b), 2), 1.5)
	n := a * math.Pow(1-e2*math.Pow(math.Sin(b), 2), -0.5)

	return ro/(m+h)*(n/a*e2*math.Sin(b)*math.Cos(b)*da+(n*n/a/a+1)*n*math.Sin(b)*math.Cos(b)*de2/2-(dx*math.Cos(l)+dy*math.Sin(l))*math.Sin(b)+dz*math.Cos(b)) - wx*math.Sin(l)*(1+e2*math.Cos(2*b)) + wy*math.Cos(l)*(1+e2*math.Cos(2*b)) - ro*ms*e2*math.Sin(b)*math.Cos(b)
}

func dL(bd, ld, h float64) float64 {
	b := bd * math.Pi / 180
	l := ld * math.Pi / 180

	n := a * math.Pow(1-e2*math.Pow(math.Sin(b), 2), -0.5)

	return ro/((n+h)*math.Cos(b))*(-dx*math.Sin(l)+dy*math.Cos(l)) + math.Tan(b)*(1-e2)*(wx*math.Cos(l)+wy*math.Sin(l)) - wz
}

func Sk42ll2Meters(lat, lon float64) (int, int, int) {
	// Номер зоны Гаусса-Крюгера
	zone := (int)(lon/6.0 + 1)

	// Параметры эллипсоида Красовского
	ka := 6378245.0                                             // Большая (экваториальная) полуось
	kb := 6356863.019                                           // Малая (полярная) полуось
	ke := (math.Pow(ka, 2) - math.Pow(kb, 2)) / math.Pow(ka, 2) // Эксцентриситет
	n := (ka - kb) / (ka + kb)                                  // Приплюснутость

	np2 := n * n
	np3 := np2 * n

	// Параметры зоны Гаусса-Крюгера
	f := 1.0                                  // Масштабный коэффициент
	lat0 := 0.0                               // Начальная параллель (в радианах)
	lon0 := float64(zone*6-3) * math.Pi / 180 // Центральный меридиан (в радианах)
	n0 := 0.0                                 // Условное северное смещение для начальной параллели
	e0 := float64(zone)*1e6 + 500000.0        // Условное восточное смещение для центрального меридиана

	// Перевод широты и долготы в радианы
	latR := lat * math.Pi / 180.0
	lonR := lon * math.Pi / 180.0

	// Вычисление переменных для преобразования
	sinLat := math.Sin(latR)
	cosLat := math.Cos(latR)
	tanLat := math.Tan(latR)

	v := ka * f * math.Pow(1-ke*math.Pow(sinLat, 2), -0.5)
	p := ka * f * (1 - ke) * math.Pow(1-ke*math.Pow(sinLat, 2), -1.5)
	n2 := v/p - 1
	m1 := (1 + n + 5.0/4.0*np2 + 5.0/4.0*np3) * (latR - lat0)
	m2 := (3*n + 3*np2 + 21.0/8.0*np3) * math.Sin(latR-lat0) * math.Cos(latR+lat0)
	m3 := (15.0/8.0*np2 + 15.0/8.0*np3) * math.Sin(2*(latR-lat0)) * math.Cos(2*(latR+lat0))
	m4 := 35.0 / 24.0 * np3 * math.Sin(3*(latR-lat0)) * math.Cos(3*(latR+lat0))
	m := kb * f * (m1 - m2 + m3 - m4)
	a1 := m + n0
	a2 := v / 2 * sinLat * cosLat
	a3 := v / 24 * sinLat * math.Pow(cosLat, 3) * (5 - math.Pow(tanLat, 2) + 9*n2)
	a3a := v / 720 * sinLat * math.Pow(cosLat, 5) * (61 - 58*math.Pow(tanLat, 2) + math.Pow(tanLat, 4))
	a4 := v * cosLat
	a5 := v / 6 * math.Pow(cosLat, 3) * (v/p - math.Pow(tanLat, 2))
	a6 := v / 120 * math.Pow(cosLat, 5) * (5 - 18*math.Pow(tanLat, 2) + math.Pow(tanLat, 4) + 14*n2 - 58*math.Pow(tanLat, 2)*n2)

	// Вычисление северного и восточного смещения (в метрах)
	nm := a1 + a2*math.Pow(lonR-lon0, 2) + a3*math.Pow(lonR-lon0, 4) + a3a*math.Pow(lonR-lon0, 6)
	em := e0 + a4*(lonR-lon0) + a5*math.Pow(lonR-lon0, 3) + a6*math.Pow(lonR-lon0, 5)

	return int(nm), int(em), zone
}

func sk42xy_to_sk42latlon(x, y int) (float64, float64) {
	// Implemented according to ГОСТ 51794 - 2001 equations: 29, 30, 31, 32, 33, 34, 35, 36
	n := float64(int(float64(y) * 0.000001))

	b := float64(x) / 6367558.4968
	b0 := b + math.Sin(2*b)*(0.00252588685-0.00001491860*(math.Pow(math.Sin(b), 2))+0.00000011904*(math.Pow(math.Sin(b), 4)))
	z0 := (float64(y) - (10*n+5)*100000) / (6378245.0 * math.Cos(b0))

	s2 := math.Pow(math.Sin(b0), 2)
	s4 := math.Pow(math.Sin(b0), 4)
	s6 := math.Pow(math.Sin(b0), 6)

	b1 := b0 - (z0*z0)*math.Sin(2*b0)*(0.251684631-0.003369263*s2+0.000011276*s4-
		(z0*z0)*(0.10500614-0.04559916*s2+0.00228901*s4-0.00002987*s6-
			(z0*z0)*(0.042858-0.025318*s2+0.014346*s4-0.001264*s6-
				(z0*z0)*(0.01672-0.00630*s2+0.01188*s4-0.00328*s6))))

	l1 := 6*(n-0.5)/57.29577951 + z0*(1-0.0033467108*s2-0.0000056002*s4-0.0000000187*s6-
		(z0*z0)*(0.16778975+0.16273586*s2-0.00052490*s4-0.00000846*s6-
			(z0*z0)*(0.0420025+0.1487407*s2-0.0059420*s4-0.0000150*s6-
				(z0*z0)*(0.01225+0.09477*s2-0.03282*s4-0.00034*s6-
					(z0*z0)*(0.0038+0.0524*s2-0.0482*s4-0.0032*s6)))))

	lon := l1 * 180 / math.Pi
	lat := b1 * 180 / math.Pi

	return lat, lon
}
