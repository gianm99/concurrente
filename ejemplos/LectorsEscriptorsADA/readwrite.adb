with Text_Io;
use  Text_Io;
with def_monitor;
use def_monitor;

procedure Readwrite is

   THREADS : constant integer := 4;
   MAX_COUNT : constant integer := 10000;
  -----
  -- Dada compartida
  -----
  counter : integer := 0;

  -----
  -- Tipus protegit per a la SC
  -----
  monitor : RWMonitor;

  -----
  -- Especificacio de la tasca
  -----
  task type Reader_writer is
    entry Start (Idx : in integer);
  end Reader_writer;

  -----
  -- Cos de la tasca
  -----
  task body Reader_writer is
    My_Idx : integer;
      c : integer;
      max : integer := MAX_COUNT / THREADS;
  begin
    accept Start (Idx : in integer) do
      My_Idx := Idx;
    end Start;
      Put_Line ("Thread " & My_Idx'img);
      for i in 0..max loop
         if ( i mod 10 = 0) then
            monitor.writerLock;
            counter := counter + 1;
            monitor.writerUnlock;
            Put_Line (My_Idx'img & " Incrementa " & counter'img);
         else
            monitor.readerLock;
            c := counter;
            monitor.readerUnlock;
            Put_Line (My_Idx'img & " Llegeix " & c'img);
         end if;

      end loop;

  end Reader_writer;

  -----
  -- Array de tasques
  -----
  type Lect_esc is array (1..THREADS) of Reader_writer;
  le : Lect_esc;

begin

  -----
  -- Start les tasques
  -----
  for Idx in 1..THREADS loop
    le(Idx).Start(Idx);
  end loop;

end Readwrite;
