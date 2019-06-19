#!/bin/sh
exec scala -savecompiled "$0" "$@"
!#

import scala.io.StdIn

import com.ibm.icu.text.RuleBasedNumberFormat
import com.ibm.icu.util.ULocale

if (args.length < 3) {
  Console.err.println("USAGE: scala batch_run_icu4j_ordinal.scala <langcode> <rulegroup> <ruleset> <numerals>")
  Console.err.println("   OR: cat <numeralfiles> | scala batch_run_icu4j_ordinal.scala <langcode> <rulegroup> <ruleset>")
  Console.err.println("   <rulegroup> cardinal/ordinal")
  System.exit(1)
}

val rulegroup = args(1)
var ruleset = args(2)
if (!ruleset.startsWith("%")) {
  ruleset = "%" + ruleset
}
val rbnf = if (rulegroup.equalsIgnoreCase("cardinal"))
  new RuleBasedNumberFormat(new ULocale(args(0)), RuleBasedNumberFormat.SPELLOUT)
else if (rulegroup.equalsIgnoreCase("ordinal"))
  new RuleBasedNumberFormat(new ULocale(args(0)), RuleBasedNumberFormat.ORDINAL)
else {
  Console.err.printf("Invalid rule group %s\n", rulegroup)
  System.exit(1)
  new RuleBasedNumberFormat(new ULocale(args(0)), RuleBasedNumberFormat.SPELLOUT)
}


if (args.length>=4) {
  for (i <- 3 until args.length) {
    val s = args(i)
    val n = if (s.matches("^[0-9-]+$")) s.toLong else s.toDouble
    val res = rbnf.format(n,ruleset)
    Console.out.println(s + "\t" + n + "\t" + res)
  }
} else {
  var s = ""
  while ({s = StdIn.readLine(); s != null}) {
    val n = if (s.matches("^[0-9-]+$")) s.toLong else s.toDouble
    val res = rbnf.format(n,ruleset)
    Console.out.println(s + "\t" + n + "\t" + res)
  }


}
