/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package getandadd;

import java.util.concurrent.atomic.AtomicInteger;

/**
 *
 * @author Daniel Muñoz i Xamena
 */

public class GetAndAdd implements Runnable {

    static final int THREADS = 2;
    static final int MAX_COUNT = 100000;
    int id;
    static AtomicInteger torn = new AtomicInteger(0);
    static int n = 0;
    static AtomicInteger number = new AtomicInteger(0);

    public GetAndAdd(int id) {
        this.id = id;
    }

    @Override
    public void run() {
        int fi = MAX_COUNT / THREADS;

        for (int i = 0; i < fi; i++) {
            AtomicInteger actual = new AtomicInteger(number.getAndAdd(1));//Incrementa la variable number i guarda el valor anterior a actual
            while (actual.get() != torn.get()) {//Espera el seu torn
            }
            n++;//Secció crítica
            torn.getAndAdd(1);//Incrementa el torn 
        }
    }

    public static void main(String[] args) throws InterruptedException {
        Thread[] threads = new Thread[THREADS];

        for (int i = 0; i < THREADS; i++) {
            threads[i] = new Thread(new GetAndAdd(i));
            threads[i].start();
        }
        for (int i = 0; i < THREADS; i++) {
            threads[i].join();
        }
        System.out.println("n = " + n);
    }

}
