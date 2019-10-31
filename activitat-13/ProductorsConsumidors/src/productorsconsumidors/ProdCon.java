/*
 * Implementau el problema de productors - consumidors usant els semàfors de 
 * Java (pot ser d'utilitat usar java.util.Deque per simular el buffer)
 */
package productorsconsumidors;

import java.util.concurrent.Semaphore;
import java.util.Deque;
import java.util.LinkedList;

/**
 *
 * @author gianm
 */
class buffer {

    //String buffer
    Deque<String> buf;

    static Semaphore mutex = new Semaphore(1);
    static Semaphore notEmpty = new Semaphore(0);

    public buffer() {
        buf = new LinkedList<String>();
    }

    public void put(int id, String data) {
        String item = String.valueOf(data);
        try {
            mutex.acquire();
            Thread.sleep(1);
        } catch (InterruptedException e) {
            System.out.println("InterruptedException caught");
        }
        buf.addLast(item);
        System.out.println("\tThread " + id + " PRODUCES: " + data);
        mutex.release();
        notEmpty.release();
    }

    void get(int id) {
        String data;
        try {
            notEmpty.acquire();
            mutex.acquire();
            Thread.sleep(3);
        } catch (InterruptedException e) {
            System.out.println("InterruptedException caught");
        }
        data = buf.pop();
        System.out.println("Thread " + id + " CONSUMES: " + data);
        mutex.release();
    }

}

class Producer implements Runnable {

    buffer buf;
    int id;

    public Producer(int id, buffer buf) {
        this.id = id;
        this.buf = buf;
    }

    @Override
    public void run() {
        System.out.println("Producer " + id);
        for (int i = 0; i < ProdCon.TO_PRODUCE; i++) {
            String data = id + " i: " + i;
            buf.put(id, data);
        }
    }

}

class Consumer implements Runnable {

    buffer buf;
    int id;

    public Consumer(int id, buffer buf) {
        this.id = id;
        this.buf = buf;
    }

    @Override
    public void run() {
        System.out.println("Consumer " + id);
        for (int i = 0; i < ProdCon.TO_PRODUCE; i++) {
            buf.get(id);
        }
    }

}

public class ProdCon {

    static final int TO_PRODUCE = 200;
    static final int PRODUCERS = 2;
    static final int CONSUMERS = 2;
    static final int THREADS = PRODUCERS + CONSUMERS;

    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws InterruptedException {
        buffer buf = new buffer();
        Thread[] threads = new Thread[THREADS];
        int i;
        // Start all threads
        for (i = 0; i < PRODUCERS; i++) {
            threads[i] = new Thread(new Producer(i, buf));
            threads[i].start();
        }
        for (i = PRODUCERS; i < THREADS; i++) {
            threads[i] = new Thread(new Consumer(i, buf));
            threads[i].start();
        }

        // Wait for all threads to complete
        for (i = 0; i < THREADS; i++) {
            threads[i].join();
        }
        System.out.println("End");
    }

}
