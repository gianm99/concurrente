package body def_monitor is
   protected body monitor is
      entry sentarse when sillas > 0 is
      begin
         sillas := sillas - 1;
         enanos_esperando:=enanos_esperando+1;
      end sentarse;

      procedure levantarse is
      begin
         sillas := sillas + 1;
      end levantarse;

      entry comer when comida > 0 is
      begin
         comida := comida - 1;
         enanos_esperando:=enanos_esperando-1;
      end comer;

      procedure darComida is
      begin
         comida := comida + 1;
      end darComida;

      procedure dormir is
      begin
         enanos_dormidos:=enanos_dormidos+1;
      end dormir;

      function esperando return Boolean is
      begin
         return enanos_esperando>0;
      end esperando;

      function dormidos return Boolean is
      begin
         return enanos_dormidos=7;
      end dormidos;

   end monitor;
end def_monitor;
