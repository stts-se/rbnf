# spellout cmd

Spell out numerals using an xml rule package from https://github.com/unicode-org/cldr/tree/master/common/rbnf

Usage:

    Usage: spellout <options> <xml file/url> <input>
      if no input argument is specified, input will be read from stdin
    Options:
      -d	Debug
      -g rule group
        	Use named rule group (default first group)
      -h	Print usage and exit
      -l	List rules and exit (rule groups and rule sets)
      -r rule set
        	Use named rule set
      -s	Check rule file syntax and exit


Example:

    $ spellout -r spellout-cardinal en.xml 97
    2019/06/11 15:27:12 Parsed rule file en.xml
    97	ninety-seven
