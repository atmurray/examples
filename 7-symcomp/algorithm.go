package main

import (
	"math"

	"gopkg.in/distil.v1"
)

//This is our distillate algorithm
type SymmetricalComponentsDistiller struct {
	// This line is required. It says this struct inherits some useful
	// default methods.
	distil.DistillateTools

	// For the frequency distillate, we make use of a rebase stage, to do
	// that, we need to know the intended frequency of the stream
	basefreq int64
}

// PadSnap is a rebase stage that will adjust the incoming data to strictly
// appear on a timebase of the given frequency (hence the 'rebase'). Any
// values that do not appear with exactly the right time are snapped to
// the nearest time, and any duplicates are dropped. In addition, any missing
// values are replaced by NaN (hence pad). The advantage of this is that it
// simplifies calculations that refer to values across time, you can rest
// assured that a value 1s ago is exactly basefreq samples away, even if
// there were holes in the data or if there were duplicates. Note that in
// the presence of duplicate data there is ZERO GUARANTEE as to WHICH of the
// multiple duplicate values you receive. In general this makes algorithms
// that compare across time quite useless, as the real time difference
// between the points a fixed interval apart will experience extreme jitter.
// The default implementation (in DistillateTools) returns RebasePassthrough
func (d *SymmetricalComponentsDistiller) Rebase() distil.Rebaser {
	return distil.RebasePadSnap(d.basefreq)
}

// This is our main algorithm. It will automatically be called with chunks
// of data that require processing by the engine.
func (d *SymmetricalComponentsDistiller) Process(in *distil.InputSet, out *distil.OutputSet) {
	/* Output 0 is real_power.
	 * Output 1 is reactive_power.
  */
	var ns int = in.NumSamples(0)
	var i int
	for i = 0; i < ns; i++ {
		var time int64 = in.Get(0, i).T

		var onetwenty = 120 * math.Pi / 180

		var magVa = in.Get(0, i).V
		var angVa = in.Get(1, i).V * math.Pi / 180
		var magVb = in.Get(2, i).V
		var angVb = in.Get(3, i).V * math.Pi / 180
                var magVc = in.Get(4, i).V
                var angVc = in.Get(5, i).V * math.Pi / 180

		var realVa = magVa * math.Cos(angVa)
		var imagVa = magVa * math.Sin(angVa)
                var realVb = magVb * math.Cos(angVb+onetwenty)
                var imagVb = magVb * math.Sin(angVb+onetwenty)
                var realVc = magVc * math.Cos(angVc-onetwenty)
                var imagVc = magVc * math.Sin(angVc-onetwenty)

		var realVp = (realVa + realVb + realVc)/3
		var imagVp = (imagVa + imagVb + imagVc)/3

                var realVd = magVb * math.Cos(angVb-onetwenty)
                var imagVd = magVb * math.Sin(angVb-onetwenty)
                var realVe = magVb * math.Cos(angVc+onetwenty)
                var imagVe = magVb * math.Sin(angVc+onetwenty)

                var realVn = (realVa + realVd + realVe)/3
                var imagVn = (imagVa + imagVd + imagVe)/3

                var realVf = magVb * math.Cos(angVb)
                var imagVf = magVb * math.Sin(angVb)
                var realVg = magVb * math.Cos(angVc)
                var imagVg = magVb * math.Sin(angVc)

                var realVz = (realVa + realVf + realVg)/3
                var imagVz = (imagVa + imagVf + imagVg)/3

		var magVp = math.Sqrt(realVp*realVp + imagVp*imagVp)
		var angVp = math.Atan(imagVp/realVp) * 180 / math.Pi
                var magVn = math.Sqrt(realVn*realVn + imagVn*imagVn)
                var angVn = math.Atan(imagVn/realVn) * 180 / math.Pi
                var magVz = math.Sqrt(realVz*realVz + imagVz*imagVz)
                var angVz = math.Atan(imagVz/realVz) * 180 / math.Pi

		if !math.IsNaN(magVp) {
			out.Add(0, time, magVp)
		}
		if !math.IsNaN(angVp) {
                        out.Add(1, time, angVp)
                }
                if !math.IsNaN(magVn) {
                        out.Add(2, time, magVn)
                }
                if !math.IsNaN(angVn) {
                        out.Add(3, time, angVn)
                }
                if !math.IsNaN(magVz) {
                        out.Add(4, time, magVz)
                }
                if !math.IsNaN(angVz) {
                        out.Add(5, time, angVz)
                }
	}
}
