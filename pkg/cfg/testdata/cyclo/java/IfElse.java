package org.example;

public class Foo {
  void CycloC(int a) {
    if (a < 0) {
      System.out.println("a");
    } else if (a == 5) {
      System.out.println("b");
    } else if (a == 6) {
      System.out.println("c");
    } else {
      System.out.println("d");
    }
    System.out.println(a);
  }
}
