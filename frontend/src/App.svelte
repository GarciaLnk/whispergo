<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Textarea } from "$lib/components/ui/textarea";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { LoadFile } from "../wailsjs/go/main/App";
  import { ProcessFile } from "../wailsjs/go/main/App";
  import ModelChooser from "./ModelChooser.svelte";

  let transcription: string = "";
  let audioFilePath: string = "";
  let audioSrc: string = "";

  const handleLoad = () => {
    LoadFile().then((response) => {
      audioFilePath = response;
      fetch(response).then((fetch_response) => {
        if (fetch_response.ok) {
          fetch_response.blob().then((blob) => {
            audioSrc = URL.createObjectURL(blob);
          });
        }
      });
    });
  };

  const handleTranscribe = () => {
    ProcessFile(audioFilePath).then((response) => {
      transcription = response;
    });
  };
</script>

<div class="container mx-auto p-4 flex flex-col h-screen">
  <div class="mb-4 flex justify-between">
    <div>
      <Button class="mb-2 w-full" on:click={handleLoad} variant="outline"
        >Load Audio File</Button
      >
      <Button
        class="btn-primary  w-full"
        on:click={handleTranscribe}
        disabled={audioFilePath === ""}>Transcribe</Button
      >
    </div>

    <div class="flex items-center ml-2">
      {#if audioSrc !== ""}
        <Label class="mr-2">{audioFilePath.split("/").pop()}</Label>
        <audio controls src={audioSrc} class="ml-auto"></audio>
      {/if}
      <div class="self-start">
        <ModelChooser />
      </div>
    </div>
  </div>

  <div class="flex-grow">
    <Textarea
      value={transcription}
      class="textarea textarea-bordered h-full"
      readonly
    />
  </div>
</div>
