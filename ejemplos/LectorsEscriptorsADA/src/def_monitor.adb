package body def_monitor is

   protected body RWMonitor is
    entry readerLock when not writing is
    begin
      readers := readers + 1;
    end readerLock;

    procedure readerUnlock is
    begin
      readers := readers - 1;
    end readerUnlock;

    entry writerLock when (readers = 0) and (not writing) is
    begin
      writing := true;
    end writerLock;

    procedure writerUnlock is
    begin
      writing := false;
    end writerUnlock;

  end RWMonitor;

end def_monitor;
