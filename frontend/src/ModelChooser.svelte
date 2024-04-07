<script lang="ts">
  import { onMount } from "svelte";
  import { Label } from "$lib/components/ui/label";
  import { Button } from "$lib/components/ui/button";
  import * as RadioGroup from "$lib/components/ui/radio-group";
  import * as Dialog from "$lib/components/ui/dialog";
  import Download from "svelte-radix/Download.svelte";
  import Gear from "svelte-radix/Gear.svelte";
  import { GetModels } from "../wailsjs/go/main/App";
  import { DownloadModel } from "../wailsjs/go/main/App";
  import { SetActiveModel } from "../wailsjs/go/main/App";
  import * as models from "../wailsjs/go/models";
  type Model = models.main.Model;

  let modelsArray = [] as Model[];
  let activeModelName = "";

  onMount(() => {
    GetModels().then((response) => {
      modelsArray = response;
      activeModelName =
        modelsArray.find((model) => model.active)?.name || modelsArray[0].name;
    });
  });

  const downloadModel = (model: Model) => {
    DownloadModel(model.name).then(() => {
      modelsList = modelsList.map((m) => {
        if (m.name === model.name) {
          m.download = true;
        }
        return m;
      });
    });
  };

  const setModel = (modelName: string) => {
    SetActiveModel(modelName);
  };

  $: modelsList = modelsArray;
</script>

<Dialog.Root>
  <Dialog.Trigger class="ml-4">
    <Gear class="h-5 w-5" />
  </Dialog.Trigger>
  <Dialog.Content class="max-h-[768px] overflow-auto">
    <Dialog.Header>
      <Dialog.Title class="mb-2">Choose a model</Dialog.Title>
      <Dialog.Description class="p-3">
        <RadioGroup.Root bind:value={activeModelName}>
          {#each modelsList as model (model.name)}
            <div class="flex items-center justify-between space-x-2 h-12">
              <div class="flex items-center space-x-2">
                <RadioGroup.Item
                  value={model.name}
                  disabled={!model.download}
                />
                <Label>{model.name} ({model.size} MB)</Label>
              </div>
              {#if !model.download}
                <Button
                  variant="outline"
                  size="icon"
                  on:click={() => downloadModel(model)}
                >
                  <Download class="h-4 w-4" />
                </Button>
              {/if}
            </div>
          {/each}
        </RadioGroup.Root>
      </Dialog.Description>
    </Dialog.Header>
    <Dialog.Footer>
      <Button type="submit" on:click={() => setModel(activeModelName)}>
        Save
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
