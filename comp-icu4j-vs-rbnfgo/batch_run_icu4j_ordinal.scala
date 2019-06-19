#!/bin/sh
exec scala -savecompiled "$0" "$@"
!#

import scala.io.StdIn

import com.ibm.icu.text.RuleBasedNumberFormat
import com.ibm.icu.util.ULocale

if (args.length < 2) {
  Console.err.println("USAGE: scala batch_run_icu4j_ordinal.scala <langcode> <ruleset> <numerals>")
  Console.err.println("   OR: cat <numeralfiles> | scala batch_run_icu4j_ordinal.scala <langcode> <ruleset>")
  System.exit(1)
}

var ruleset = args(1)
if (!ruleset.startsWith("%")) {
  ruleset = "%" + ruleset
}
val rbnf = new RuleBasedNumberFormat(new ULocale(args(0)), RuleBasedNumberFormat.SPELLOUT)

if (args.length>=3) {
  for (i <- 2 until args.length) {
    val s = args(i)
    val n = s.toLong
    val res = rbnf.format(n,ruleset)
    Console.out.println(s + "\t" + res)
  }
} else {
  var s = ""
  while ({s = StdIn.readLine(); s != null}) {
    val n = s.toLong
    val res = rbnf.format(n,ruleset)
    Console.out.println(s + "\t" + res)
  }


}
