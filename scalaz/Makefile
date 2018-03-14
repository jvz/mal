TESTS =

SOURCES_BASE = src/main/scala/types.scala src/main/scala/reader.scala src/main/scala/printer.scala
SOURCES_LISP = src/main/scala/env.scala src/main/scala/core.scala src/main/scala/stepA_mal.scala
SOURCES = $(SOURCES_BASE) $(SOURCES_LISP)

TARGET_DIR=target/scala-2.12

all: $(TARGET_DIR)/mal.jar

dist: mal

mal: $(TARGET_DIR)/mal.jar
	cp $< $@

$(TARGET_DIR)/mal.jar:
	./sbt assembly

$(TARGET_DIR)/classes/step%.class: step%.scala $(SOURCES)
	./sbt assembly

clean:
	rm -rf mal target

.PHONY: all dist clean stats tests $(TESTS)

stats: $(SOURCES)
	@wc $^
	@printf "%5s %5s %5s %s\n" `grep -E "^[[:space:]]*//|^[[:space:]]*$$" $^ | wc` "[comments/blanks]"
stats-lisp: $(SOURCES_LISP)
	@wc $^
	@printf "%5s %5s %5s %s\n" `grep -E "^[[:space:]]*//|^[[:space:]]*$$" $^ | wc` "[comments/blanks]"

tests: $(TESTS)

$(TESTS):
	@echo "Running $@"; \
	ruby $@ || exit 1; \
