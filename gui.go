//heidi hoehn
//Zweck: Zahlen mit der Maus zeichnen und in einem 2dimensionalen Array speichern
//Datum: 30.03.2022

package gui

import (
	"fmt"
	"gfx"
	"time"
	"gonum.org/v1/gonum/mat"
)

const vergroesserung uint16 = 20
const SCHWARZ uint16 = 0
const WEISS uint16 = 255
const GRAU uint16 = 128
// Vor.: Grafikfenster muss offen sein.
// Eff.: Benutzer zeichnet mit der Maus eine Zahl, indem er die linke Maustaste gedrückt hält und die Maus bewegt
// 		 Wird die rechte Maustaste gedrückt, signalisiert der Benutzer: "meine Zeichung ist fertig!"
// Erg.: 2D-Array, das das gezeichnete Bild mit in Schwarz(= 0)-Weiß(= 255)- Grau(= 128)kodiert

func EinlesenZeichnung() [560][560]uint16 {
	var erg [560][560]uint16 //Variable für das Ergebnis erstellen
	for i := 0; i < len(erg); i++ {
		for j := 0; j < len(erg); j++ {
				erg[i][j] = WEISS
		}
	}
	gfx.Stiftfarbe(0, 0, 0)  // setzt Stiftfarbe auf schwarz
	gfx.Schreibe(0, 0, "Zeichne bitte eine Zahl!")
	gfx.Schreibe(0, 12, "Wenn du fertig bist, drueck die rechte Maustaste!")
	for { //Endlosschleife, in der die Benutzereingabe abgefragt wird
		taste, status, mausX, mausY := gfx.MausLesen1() //Auslesen der Mauseigenschaften
		if taste == 1 && status != -1 {                 // Wenn Maustaste 1 gedrückt oder gehalten wird, ...
			gfx.Vollrechteck(mausX, mausY, vergroesserung, vergroesserung) //zeichnen wir einen "Pixel" (= Vollrechteck) der Größe "vergroesserung"
			var i, j uint16
			for i = 0; i < vergroesserung; i++ {
				for j = 0; j < vergroesserung; j++ {
					erg[mausX+i][mausY+j] = SCHWARZ
				}
			}
		}
		if taste == 3 {
		//graue Umrandung für die schwarzen Pixel zeichnen; Grundidee: zeichne um jeden schwarzen Pixel herum: graue Pixel - nachdem Benutzer Zahl gemalt hat, keinesfalls: schwarze Pixel überschreiben
		var i, j, di, dj int		
			for i = 1; i < len(erg)-1; i++ {
				for j = 1; j < len(erg)-1; j++ {
					
					if erg[i][j] == SCHWARZ {
						for di = -1; di < 2; di++ {
							for dj = -1; dj < 2; dj++ {
								if erg[i+di][j+dj] != SCHWARZ{
									gfx.Stiftfarbe(128, 128, 128)  
									gfx.Vollrechteck(uint16(i+di),uint16(j+dj), 10, 10)
									erg[i+di][j+dj] = GRAU								
									}
							}	
						}
						
					}
				}
			}	
			
			
			time.Sleep(4*time.Second)
			break
		}
	}
	return erg
}

//
//			gfx.Vollrechteck(mausX-5, mausY-5, vergroesserung+10, vergroesserung+10) //zeichnen wir einen "Pixel" (= Vollrechteck) der Größe "vergroesserung"
//			gfx.Stiftfarbe(0, 0, 0)  // setzt Stiftfarbe auf schwarz

func bildSkalieren(bildOriginal [560][560]uint16) [28][28]float64 { // Runterskalieren des Bildes von 560x560 auf 28x28
	var bildScaled [28][28]float64
	row := 0
	for x := 0; x < 560; {

		col := 0
		for i := 0; i < 560; {

			summe := 0
			for y := x; y < x+20; y++ {
				for j := i; j < i+20; j++ {
					summe = summe + int(bildOriginal[j][y]) //Bilder sind spaltenweise angeordnet, nicht zeilenweise!
				}
			}

			summeScaled := (float64(summe) / (20.0 * 20.0* 255))
			bildScaled[row][col] = summeScaled

			i = i + 20
			col++
		}
		x = x + 20
		row++
	}
	return bildScaled
}

func gibBildMatrix(bildArray [28][28]float64) *mat.Dense {
	bildSlice := make([]float64, 28*28)

	index := 0
	for i := 0; i < 28; i++ {
		for j := 0; j < 28; j++ {
			bildSlice[index] = bildArray[i][j]
			index++
		}
	}
	bildMatrix := mat.NewDense(784, 1, bildSlice)

	return bildMatrix

}

//Vor.: keine
//Eff.: ein Grafikfenster mit weißem Hintergrund wird erstellt, die Funktion "EinlesenZeichnung()" wird aufgerufen, damit Benutzer eine Zahl zeichnen kann.
//		Wenn der Benutzer fertig ist, wird zuerst die Kodierung des Bildes im Terminal ausgegeben und dann aus der Kodierung wieder das gezeichnete Bild rekonstruiert.
//		Der letzte Schritt stellt sicher, dass die Kodierung korrekt ist.
func ZahlMalen() *mat.Dense {
	gfx.Fenster(560, 560) // öffnet das Grafikfenster mit weißem Hintergrund
	fmt.Println(gfx.GibFont())

	var bildBinaer [560][560]uint16
	bildBinaer = EinlesenZeichnung()
	//fmt.Println(bildBinaer) //gib die Kodierung im Terminal aus
	// rekonstruiere aus der Kodierung das gezeichnete Bild
	gfx.Stiftfarbe(255, 255, 255) // weiß
	gfx.Cls()
	gfx.Stiftfarbe(0, 0, 0) // schwarz
	bildScale := bildSkalieren(bildBinaer)
	for i := 0; i < 28; i++ {
		for j := 0; j < 28; j++ {
			if bildScale[j][i] < 0.5 { //Durchschnittswert kleiner 0.5 bedeutet "schwarz"; Hinweis: j und i getauscht wegen Spiegelung schneller Workaround
				gfx.Punkt(uint16(i), uint16(j))
			}
		}
	}
	gfx.TastaturLesen1()
	zahlMatrix := gibBildMatrix(bildScale)

	return zahlMatrix

}
