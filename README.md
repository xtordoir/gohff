# gohff

A library implementing basis blocks for coastline and overshoot trading as described in:

http://www.olsendata.com/fileadmin/Publications/Client_Papers/201206-DupuisOlsen-HFFinScaLawTraMod.pdf

`hff.go`

**Overshoot** gives the market direction and extension for a given scale.


`spectrum.go`

**Spectrum** Overshoots for a range of scales

`trader.go`

**TraderProcess** holds the state of a trading agent. It is used to know current exposure of the agent and accumulated PL from history of increase and decrease of positions.

`trader/coastline.go`

**CoastlineTrade** is an agent for the Coastline trade. It uses **Parameters** (initialization params and current agent state) and **TraderProcess** to track exposure and PL.

 
