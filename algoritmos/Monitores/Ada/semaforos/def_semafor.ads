protected type Semafor(Inicial: Natural) is
	entry Wait;
	procedure Signal;
private
	Contador: Natural := Inicial;
end Semafor;