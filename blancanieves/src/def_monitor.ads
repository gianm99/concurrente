-- Autor: Gian Lucas Martín Chamorro
-- Enlace: https://youtu.be/7UQQAyvB-dE
package def_monitor is

   protected type monitor is
      entry sentarse;
      procedure levantarse;
      entry comer;
      procedure darComida;
      procedure dormir;
      function esperando return integer;
      function dormidos return integer;
      function sillas_libres return integer;
   private
      -- Cantidad de sillas libres
      sillas: integer := 4;
      -- Cantidad de enanos esperando para comer
      enanos_esperando: integer := 0;
      -- Cantidad de enanos dormidos
      enanos_dormidos: integer := 0;
      -- Cantidad de platos de comida listos para comer
      comida: integer := 0;
   end monitor;

end def_monitor;
