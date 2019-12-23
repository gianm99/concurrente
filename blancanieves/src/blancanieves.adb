with Ada.Text_IO,Ada.Integer_Text_IO,def_monitor,Ada.Strings.Unbounded;
use Ada.Text_IO,Ada.Integer_Text_IO,def_monitor,Ada.Strings.Unbounded;

procedure blancanieves is
   type nombres_enanos is(SABIO,GRUNON,FELIZ,DORMILON,TIMIDO,MOCOSO,TONTIN);
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
      max: integer := 2;
   begin
      accept Start (nombre: in nombres_enanos) do
         mi_nombre:=nombre;
      end Start;
      for I in 1..max loop
         -- VA A TRABAJAR
         Put_Line(mi_nombre'img & " va a trabajar a la mina");
         delay Duration(4.0);
         -- SE SIENTA EN UNA SILLA
         m.sentarse;
         Put_Line(mi_nombre'img & " se sienta");
         -- COME
         m.comer;
         Put_Line("-----------------> "& mi_nombre'img & " come!!!");
         Put_Line("esperando = " & m.esperando'img & " sillas = " & m.sillas_libres'img);
         delay Duration(1.5);
         -- SE LEVANTA DE LA SILLA
         m.levantarse;
         Put_Line("esperando = " & m.esperando'img & " sillas = " & m.sillas_libres'img);
      end loop;
      -- SE VA A DORMIR
      m.dormir;
      Put_Line(mi_nombre'img & " se va a DORMIR ");
   end tarea_enano;

   task body tarea_blancanieves is
   begin
      accept Start;
      while m.dormidos<7 loop
         -- COMPRUEBA SI HAY ENANOS ESPERANDO PARA COMER
         while m.esperando>0 loop
            -- MIENTRAS HAYA ENANOS ESPERANDO REPARTE COMIDA
            Put_Line("BLANCANIEVES cocina para un enano");
            delay Duration(0.5);
            m.darComida;
         end loop;
         -- SALE A PASEAR
         Put_Line("BLANCANIEVES se va a pasear <-----------------");
         delay Duration(1.5);
      end loop;
      -- SE VA A DORMIR
      Put_Line("BLANCANIEVES se va a dormir");
   end tarea_blancanieves;

   -----------------------------
   --  Array de enanos
   -----------------------------
   type Enanos is array (nombres_enanos) of tarea_enano;
   e:Enanos;
   type Blancanieves is new tarea_blancanieves;
   b: Blancanieves;
begin
   b.Start; -- Empieza Blancanieves
   for I in e'Range loop
      e(I).Start(I);
   end loop;
end blancanieves;
