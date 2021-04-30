package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/runners/direct"
)

func main() {
	// In order to start creating the pipeline for execution, a Pipeline object is needed.
	p := beam.NewPipeline()
	s := p.Root()

	// The pipeline object encapsulates all the data and steps in your processing task.
	// It is the basis for creating the pipeline's data sets as PCollections and its operations
	// as transforms.

	// The PCollection abstraction represents a potentially distributed,
	// multi-element data set. You can think of a PCollection as “pipeline” data;
	// Beam transforms use PCollection objects as inputs and outputs. As such, if
	// you want to work with data in your pipeline, it must be in the form of a
	// PCollection.

	// Transformations are applied in a scoped fashion to the pipeline. The scope
	// can be obtained from the pipeline object.

	// Start by reading text from an input files, and receiving a PCollection.
	// lines := textio.Read(s, "protocol://path/file*.txt")
	lines := textio.Read(s, "hamlet.txt")

	// Transforms are added to the pipeline so they are part of the work to be
	// executed.  Since this transform has no PCollection as an input, it is
	// considered a 'root transform'

	// A pipeline can have multiple root transforms
	moreLines := textio.Read(s, "julius.txt")

	// Further transforms can be applied, creating an arbitrary, acyclic graph.
	// Subsequent transforms (and the intermediate PCollections they produce) are
	// attached to the same pipeline.
	all := beam.Flatten(s, lines, moreLines)
	wordRegexp := regexp.MustCompile(`[a-zA-Z]+('[a-z])?`)
	words := beam.ParDo(s, func(line string, emit func(string)) {
		for _, word := range wordRegexp.FindAllString(line, -1) {
			emit(word)
		}
	}, all)
	formatted := beam.ParDo(s, strings.ToUpper, words)
	textio.Write(s, "./wordcounts.txt", formatted)

	// Applying a transform adds it to the pipeline, rather than executing it
	// immediately.  Once the whole pipeline of transforms is constructed, the
	// pipeline can be executed by a PipelineRunner.  The direct runner executes the
	// transforms directly, sequentially, in this one process, which is useful for
	// unit tests and simple experiments:
	_, err := direct.Execute(context.Background(), p)
	if err != nil {
		fmt.Printf("Pipeline failed: %v", err)
	}
}
