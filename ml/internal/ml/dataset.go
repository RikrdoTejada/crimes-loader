// Package ml implementa un Random Forest en Go puro con entrenamiento
// paralelo mediante goroutines y channels (patrón productor-consumidor).
package ml

import (
	"math/rand"
	"sort"

	"ml-paralelo/internal/models"
)

// NumFeatures es el número de features de entrada del modelo.
// Orden: hour, district, primary_type (encoded), domestic,
// day_of_week, month, community_area, beat
const NumFeatures = 8

// Nombres de las features, en el mismo orden que el vector de entrada.
var FeatureNames = [NumFeatures]string{
	"hour", "district", "primary_type", "domestic",
	"day_of_week", "month", "community_area", "beat",
}

// Dataset contiene las muestras codificadas listas para entrenamiento.
// X[i] es el vector de features de la muestra i; Y[i] es su etiqueta
// (1 = arrest, 0 = no arrest).
type Dataset struct {
	X [][NumFeatures]float64
	Y []int
	// TypeEncoder mapea cada primary_type a su código entero (label encoding).
	TypeEncoder map[string]int
}

// BuildDataset convierte los CleanRecord en un Dataset numérico.
// Realiza el label encoding de primary_type en una sola pasada,
// asignando códigos en orden alfabético para que el encoding sea
// determinista entre ejecuciones.
func BuildDataset(records []models.CleanRecord) *Dataset {
	// 1. Recolectar los tipos únicos y ordenarlos alfabéticamente
	typeSet := make(map[string]struct{})
	for i := range records {
		typeSet[records[i].PrimaryType] = struct{}{}
	}
	types := make([]string, 0, len(typeSet))
	for t := range typeSet {
		types = append(types, t)
	}
	sort.Strings(types)

	encoder := make(map[string]int, len(types))
	for i, t := range types {
		encoder[t] = i
	}

	// 2. Construir las matrices X e Y
	ds := &Dataset{
		X:           make([][NumFeatures]float64, len(records)),
		Y:           make([]int, len(records)),
		TypeEncoder: encoder,
	}
	for i := range records {
		r := &records[i]
		domestic := 0.0
		if r.Domestic {
			domestic = 1.0
		}
		ds.X[i] = [NumFeatures]float64{
			float64(r.Hour),
			float64(r.District),
			float64(encoder[r.PrimaryType]),
			domestic,
			float64(r.DayOfWeek),
			float64(r.Month),
			float64(r.CommunityArea),
			float64(r.Beat),
		}
		if r.Arrest {
			ds.Y[i] = 1
		}
	}
	return ds
}

// StratifiedSplit divide el dataset en entrenamiento y validación
// preservando la proporción de clases (estratificación por Y).
// trainFrac es la fracción destinada a entrenamiento (ej. 0.8).
func (ds *Dataset) StratifiedSplit(trainFrac float64, rng *rand.Rand) (train, test *Dataset) {
	// separar índices por clase
	var pos, neg []int
	for i, y := range ds.Y {
		if y == 1 {
			pos = append(pos, i)
		} else {
			neg = append(neg, i)
		}
	}
	rng.Shuffle(len(pos), func(i, j int) { pos[i], pos[j] = pos[j], pos[i] })
	rng.Shuffle(len(neg), func(i, j int) { neg[i], neg[j] = neg[j], neg[i] })

	nPosTrain := int(float64(len(pos)) * trainFrac)
	nNegTrain := int(float64(len(neg)) * trainFrac)

	trainIdx := append(append([]int{}, pos[:nPosTrain]...), neg[:nNegTrain]...)
	testIdx := append(append([]int{}, pos[nPosTrain:]...), neg[nNegTrain:]...)

	build := func(idx []int) *Dataset {
		sub := &Dataset{
			X:           make([][NumFeatures]float64, len(idx)),
			Y:           make([]int, len(idx)),
			TypeEncoder: ds.TypeEncoder,
		}
		for i, j := range idx {
			sub.X[i] = ds.X[j]
			sub.Y[i] = ds.Y[j]
		}
		return sub
	}
	return build(trainIdx), build(testIdx)
}

// ClassWeights calcula los pesos balanceados de cada clase:
// peso_c = total / (2 * count_c). La clase minoritaria recibe
// mayor peso, compensando el desbalance ~87/13 del dataset.
func (ds *Dataset) ClassWeights() (w0, w1 float64) {
	var n0, n1 int
	for _, y := range ds.Y {
		if y == 1 {
			n1++
		} else {
			n0++
		}
	}
	total := float64(n0 + n1)
	if n0 > 0 {
		w0 = total / (2.0 * float64(n0))
	}
	if n1 > 0 {
		w1 = total / (2.0 * float64(n1))
	}
	return w0, w1
}
