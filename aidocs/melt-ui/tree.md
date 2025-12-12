# Tree

A hierarchical list component for displaying nested data with expand/collapse functionality.

## Use Cases

- File explorers
- Navigation menus
- Category trees
- Organization charts

## Installation

```typescript
import { createTreeView } from '@melt-ui/svelte';
```

## API Reference

### createTreeView Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultExpanded` | `string[]` | `[]` | Initially expanded items |
| `expanded` | `Writable<string[]>` | — | Controlled expanded items |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `onExpandedChange` | `ChangeFn` | — | Expanded change callback |

### For Selection (Optional)

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultSelectedId` | `string` | — | Initially selected item |
| `selectedId` | `Writable<string>` | — | Controlled selection |
| `onSelectedChange` | `ChangeFn` | — | Selection callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `tree` | Root tree container |
| `item` | Tree item container |
| `group` | Nested group container |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `expanded` | `Writable<string[]>` | Expanded item IDs |
| `selectedId` | `Writable<string>` | Selected item ID |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isExpanded(id)` | Check if item is expanded |
| `isSelected(id)` | Check if item is selected |

## Data Attributes

### Tree
- `[data-melt-tree]` - Present on tree root

### Item
- `[data-id]` - Item identifier
- `[data-value]` - Item value
- `[data-melt-tree-item]` - Present on items

### Group
- `[data-melt-tree-group]` - Present on groups

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Down` | Next visible item |
| `Arrow Up` | Previous visible item |
| `Arrow Right` | Expand item or move to first child |
| `Arrow Left` | Collapse item or move to parent |
| `Home` | First item |
| `End` | Last visible item |
| `Enter` / `Space` | Select item |

## Accessibility

Follows [WAI-ARIA Tree Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/treeview/):
- `role="tree"` and `role="treeitem"`
- `aria-expanded` for expandable items
- Full keyboard navigation

## Example

```svelte
<script lang="ts">
  import { createTreeView, melt } from '@melt-ui/svelte';
  import { ChevronRight, Folder, File } from 'lucide-svelte';

  type TreeItem = {
    id: string;
    name: string;
    children?: TreeItem[];
  };

  const data: TreeItem[] = [
    {
      id: 'src',
      name: 'src',
      children: [
        { id: 'app', name: 'app.ts' },
        {
          id: 'components',
          name: 'components',
          children: [
            { id: 'button', name: 'Button.svelte' },
            { id: 'input', name: 'Input.svelte' }
          ]
        }
      ]
    },
    { id: 'readme', name: 'README.md' },
    { id: 'package', name: 'package.json' }
  ];

  const {
    elements: { tree, item, group },
    states: { expanded },
    helpers: { isExpanded, isSelected }
  } = createTreeView({
    defaultExpanded: ['src']
  });
</script>

{#snippet treeItem(node: TreeItem, level: number)}
  <div
    use:melt={$item({ id: node.id, hasChildren: !!node.children?.length })}
    class="flex items-center gap-2 px-2 py-1 rounded cursor-pointer hover:bg-ocean-100
      data-[selected]:bg-ocean-500 data-[selected]:text-white"
    style="padding-left: {level * 16 + 8}px"
  >
    {#if node.children?.length}
      <ChevronRight
        class="h-4 w-4 transition-transform {$isExpanded(node.id) ? 'rotate-90' : ''}"
      />
      <Folder class="h-4 w-4 text-ocean-500" />
    {:else}
      <span class="w-4"></span>
      <File class="h-4 w-4 text-ocean-400" />
    {/if}
    <span class="text-sm">{node.name}</span>
  </div>

  {#if node.children?.length && $isExpanded(node.id)}
    <div use:melt={$group({ id: node.id })}>
      {#each node.children as child}
        {@render treeItem(child, level + 1)}
      {/each}
    </div>
  {/if}
{/snippet}

<div use:melt={$tree} class="w-64 p-2 bg-white rounded-lg border border-ocean-200">
  {#each data as node}
    {@render treeItem(node, 0)}
  {/each}
</div>
```

## Styling with Tailwind

```css
/* Tree container */
[data-melt-tree] {
  @apply space-y-0.5;
}

/* Tree item */
[data-melt-tree-item] {
  @apply flex items-center gap-2 px-2 py-1 rounded cursor-pointer;
  @apply hover:bg-ocean-100 transition-colors;
}

/* Selected item */
[data-melt-tree-item][data-selected] {
  @apply bg-ocean-500 text-white;
}

/* Focus state */
[data-melt-tree-item]:focus-visible {
  @apply ring-2 ring-ocean-500 ring-inset;
}

/* Group indentation */
[data-melt-tree-group] {
  @apply ml-4;
}
```

## Navigation Tree

```svelte
<script lang="ts">
  import { createTreeView, melt } from '@melt-ui/svelte';
  import { ChevronRight } from 'lucide-svelte';

  const {
    elements: { tree, item, group },
    helpers: { isExpanded }
  } = createTreeView();

  const navigation = [
    {
      id: 'dashboard',
      label: 'Dashboard',
      href: '/dashboard'
    },
    {
      id: 'team',
      label: 'Team',
      children: [
        { id: 'calendar', label: 'Calendar', href: '/team/calendar' },
        { id: 'requests', label: 'Requests', href: '/team/requests' }
      ]
    },
    {
      id: 'admin',
      label: 'Admin',
      children: [
        { id: 'users', label: 'Users', href: '/admin/users' },
        { id: 'settings', label: 'Settings', href: '/admin/settings' }
      ]
    }
  ];
</script>

<nav use:melt={$tree} class="w-56 space-y-1">
  {#each navigation as navItem}
    {#if navItem.children}
      <div>
        <button
          use:melt={$item({ id: navItem.id, hasChildren: true })}
          class="w-full flex items-center justify-between px-3 py-2 rounded-lg text-ocean-700 hover:bg-ocean-100"
        >
          <span>{navItem.label}</span>
          <ChevronRight
            class="h-4 w-4 transition-transform {$isExpanded(navItem.id) ? 'rotate-90' : ''}"
          />
        </button>

        {#if $isExpanded(navItem.id)}
          <div use:melt={$group({ id: navItem.id })} class="ml-4 mt-1 space-y-1">
            {#each navItem.children as child}
              <a
                use:melt={$item({ id: child.id })}
                href={child.href}
                class="block px-3 py-2 rounded-lg text-sm text-ocean-600 hover:bg-ocean-100
                  data-[selected]:bg-ocean-500 data-[selected]:text-white"
              >
                {child.label}
              </a>
            {/each}
          </div>
        {/if}
      </div>
    {:else}
      <a
        use:melt={$item({ id: navItem.id })}
        href={navItem.href}
        class="block px-3 py-2 rounded-lg text-ocean-700 hover:bg-ocean-100
          data-[selected]:bg-ocean-500 data-[selected]:text-white"
      >
        {navItem.label}
      </a>
    {/if}
  {/each}
</nav>
```

## With Checkboxes

```svelte
<script lang="ts">
  import { createTreeView, melt } from '@melt-ui/svelte';
  import { Check } from 'lucide-svelte';

  let checkedItems = $state<Set<string>>(new Set());

  function toggleCheck(id: string) {
    if (checkedItems.has(id)) {
      checkedItems.delete(id);
    } else {
      checkedItems.add(id);
    }
    checkedItems = new Set(checkedItems);
  }
</script>

{#snippet treeItem(node)}
  <div
    use:melt={$item({ id: node.id, hasChildren: !!node.children })}
    class="flex items-center gap-2 px-2 py-1"
  >
    <button
      onclick={() => toggleCheck(node.id)}
      class="h-4 w-4 rounded border border-ocean-400 flex items-center justify-center
        {checkedItems.has(node.id) ? 'bg-ocean-500 border-ocean-500' : ''}"
    >
      {#if checkedItems.has(node.id)}
        <Check class="h-3 w-3 text-white" />
      {/if}
    </button>
    <span>{node.name}</span>
  </div>
{/snippet}
```

## Controlled Expansion

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';

  const expandedItems = writable(['root', 'src']);

  const tree = createTreeView({
    expanded: expandedItems,
    onExpandedChange: ({ next }) => {
      console.log('Expanded items:', next);
      return next;
    }
  });

  // Programmatic control
  function expandAll() {
    expandedItems.set(['root', 'src', 'components', 'lib']);
  }

  function collapseAll() {
    expandedItems.set([]);
  }
</script>
```

## Async Loading

```svelte
<script lang="ts">
  let loadedChildren = $state<Map<string, TreeItem[]>>(new Map());
  let loadingIds = $state<Set<string>>(new Set());

  async function loadChildren(id: string) {
    if (loadedChildren.has(id)) return;

    loadingIds.add(id);
    loadingIds = new Set(loadingIds);

    const children = await fetchChildren(id);

    loadedChildren.set(id, children);
    loadedChildren = new Map(loadedChildren);

    loadingIds.delete(id);
    loadingIds = new Set(loadingIds);
  }
</script>

{#snippet treeItem(node)}
  <div
    use:melt={$item({ id: node.id, hasChildren: node.hasChildren })}
    onclick={() => node.hasChildren && loadChildren(node.id)}
  >
    {#if loadingIds.has(node.id)}
      <Loader class="h-4 w-4 animate-spin" />
    {:else if node.hasChildren}
      <ChevronRight class="h-4 w-4 {$isExpanded(node.id) ? 'rotate-90' : ''}" />
    {/if}
    <span>{node.name}</span>
  </div>

  {#if $isExpanded(node.id) && loadedChildren.has(node.id)}
    <div use:melt={$group({ id: node.id })}>
      {#each loadedChildren.get(node.id) as child}
        {@render treeItem(child)}
      {/each}
    </div>
  {/if}
{/snippet}
```
