with Ada.Text_IO,Ada.Integer_Text_IO,def_monitor,Ada.Strings.Unbounded;
use Ada.Text_IO,Ada.Integer_Text_IO,def_monitor,Ada.Strings.Unbounded;

procedure blancanieves is
   type nombres_enanos is(SABIO,GRUNON,FELIZ,DORMILON,TIMIDO,MOCOSO,TONTIN);
   MAX: constant integer:=2;
   -----------------------------
   --  Tipo protegido para la SC
   -----------------------------
   m: monitor;
   -----------------------------
   --  Especificaciones de las tareas
   -----------------------------
   task type tarea_enano is
      entry Start(nombre: in nombres_enanos);
   end tarea_enano;

   task type tarea_blancanieves is
      entry Start;
   end tarea_blancanieves;

   -----------------------------
   -- Cuerpos de las tareas
   -----------------------------
   task body tarea_enano is
      mi_nombre : nombres_enanos;
   begin
      accept Start (nombre: in nombres_enanos) do
         mi_nombre:=nombre;
      end Start;
      for I in 1..MAX loop
         Put_Line(mi_nombre'img & " va a trabajar a la mina");
         delay Duration(4.0);
         m.sentarse;
         Put_Line(mi_nombre'img & " se sienta");
         m.comer;
         Put_Line("-----------------> "& mi_nombre'img & " come!!!");
         Put_Line("esperando = " & m.esperando'img & " sillas = " & m.sillas_libres'img);
         delay Duration(1.5);
         m.levantarse;
         Put_Line("esperando = " & m.esperando'img & " sillas = " & m.sillas_libres'img);
      end loop;
      m.dormir;
      Put_Line(mi_nombre'img & " se va a DORMIR ");
   end tarea_enano;

   task body tarea_blancanieves is
   begin
      accept Start;
      while m.dormidos<7 loop
         while m.esperando>0 loop
            Put_Line("BLANCANIEVES cocina para un enano");
            delay Duration(0.5);
            m.darComida;
         end loop;
         Put_Line("BLANCANIEVES se va a pasear <-----------------");
         delay Duration(1.5);
      end loop;
      Put_Line("BLANCANIEVES se va a dormir");
   end tarea_blancanieves;

   -----------------------------
   --  Tareas
   -----------------------------
   type Enanos is array (nombres_enanos) of tarea_enano;
   e:Enanos;
   type Blancanieves is new tarea_blancanieves;
   b: Blancanieves;
begin
   -- Empezar todas las tareas
   b.Start; -- Blancanieves
   for I in e'Range loop
      e(I).Start(I); -- Enanos
   end loop;
end blancanieves;
