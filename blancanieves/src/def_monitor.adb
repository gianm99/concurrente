-- Autor: Gian Lucas Martín Chamorro
-- Enlace: https://youtu.be/7UQQAyvB-dE
package body def_monitor is
   protected body monitor is
      entry sentarse when sillas > 0 is
      begin
         -- Se ocupa una silla
         sillas := sillas - 1;
         enanos_esperando:=enanos_esperando+1;
      end sentarse;

      procedure levantarse is
      begin
         -- Se libera una silla
         sillas := sillas + 1;
      end levantarse;

      entry comer when comida > 0 is
      begin
         -- Disminuye en uno la cantidad de comida
         comida := comida - 1;
         -- Hay un enano menos esperando para comer
         enanos_esperando:=enanos_esperando-1;
      end comer;

      procedure darComida is
      begin
         -- Blancanieves ha cocinado un plato de comida
         comida := comida + 1;
      end darComida;

      procedure dormir is
      begin
         -- Aumenta en uno la cantidad de enanos dormidos
         enanos_dormidos:=enanos_dormidos+1;
      end dormir;

      function esperando return integer is
      begin
         -- Devuelve la cantidad de enanos esperando para comer
         return enanos_esperando;
      end esperando;

      function dormidos return integer is
      begin
         -- Devuelve la cantidad de enanos dormidos
         return enanos_dormidos;
      end dormidos;

      function sillas_libres return integer is
      begin
         -- Devuelve la cantidad de sillas libres
         return sillas;
      end sillas_libres;
   end monitor;
end def_monitor;
