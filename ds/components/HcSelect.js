import {
  computed,
  defineComponent,
  h,
  nextTick,
  onBeforeUnmount,
  onMounted,
  ref,
  watch,
} from 'vue';

let selectIdCounter = 0;

export default defineComponent({
  name: 'HcSelect',
  props: {
    id: {
      type: String,
      default: null,
    },
    modelValue: {
      type: [String, Number, Boolean, Object],
      default: null,
    },
    options: {
      type: Array,
      default: () => [],
    },
    placeholder: {
      type: String,
      default: 'Selecciona una opción',
    },
    disabled: Boolean,
  },
  emits: ['update:modelValue', 'change', 'manage-projects'],
  setup(props, { emit }) {
    const isOpen = ref(false);
    const triggerRef = ref(null);
    const listRef = ref(null);
    const searchInputRef = ref(null);
    const activeIndex = ref(-1);
    const searchQuery = ref('');

    const componentId = ref(props.id ?? `hc-select-${++selectIdCounter}`);
    watch(
      () => props.id,
      (value) => {
        if (value) {
          componentId.value = value;
        }
      },
    );

    const listboxId = computed(() => `${componentId.value}-listbox`);
    const optionId = (index) => `${componentId.value}-option-${index}`;

    const selectedOption = computed(
      () => props.options.find((option) => option.value === props.modelValue) ?? null,
    );

    const shouldShowSearch = computed(() => props.options.length > 6);

    const normalizedQuery = computed(() => searchQuery.value.trim().toLowerCase());

    const visibleOptions = computed(() => {
      if (!shouldShowSearch.value) {
        return props.options;
      }
      const query = normalizedQuery.value;
      if (!query) {
        return props.options;
      }
      return props.options.filter((option) => {
        const label = typeof option.label === 'string' ? option.label.toLowerCase() : '';
        const description =
          typeof option.description === 'string' ? option.description.toLowerCase() : '';
        const meta = typeof option.meta === 'string' ? option.meta.toLowerCase() : '';
        return label.includes(query) || description.includes(query) || meta.includes(query);
      });
    });

    const setActiveIndex = (index) => {
      activeIndex.value = index;
      nextTick(() => {
        if (index < 0) return;
        const listEl = listRef.value;
        if (!listEl) return;
        const optionEl = listEl.children?.[index];
        if (optionEl && optionEl.scrollIntoView) {
          optionEl.scrollIntoView({ block: 'nearest' });
        }
      });
    };

    const setActiveByValue = (value) => {
      if (value === null || value === undefined) {
        setActiveIndex(-1);
        return;
      }
      const options = visibleOptions.value;
      const idx = options.findIndex((option) => option.value === value);
      setActiveIndex(idx);
    };

    watch(
      () => props.modelValue,
      (value) => {
        setActiveByValue(value);
      },
    );

    const openDropdown = () => {
      if (props.disabled || !props.options.length) return;
      if (isOpen.value) return;
      searchQuery.value = '';
      isOpen.value = true;
      nextTick(() => {
        if (selectedOption.value) {
          setActiveByValue(selectedOption.value.value);
        } else {
          setActiveIndex(visibleOptions.value.length ? 0 : -1);
        }
        if (shouldShowSearch.value) {
          searchInputRef.value?.focus({ preventScroll: true });
        } else {
          listRef.value?.focus({ preventScroll: true });
        }
      });
    };

    const closeDropdown = (restoreFocus = true) => {
      if (!isOpen.value) return;
      isOpen.value = false;
      setActiveIndex(-1);
      searchQuery.value = '';
      if (restoreFocus) {
        nextTick(() => {
          triggerRef.value?.focus({ preventScroll: true });
        });
      }
    };

    const toggleDropdown = () => {
      if (isOpen.value) {
        closeDropdown();
      } else {
        openDropdown();
      }
    };

    const selectOption = (option) => {
      if (!option) return;
      emit('update:modelValue', option.value);
      emit('change', option.value);
      closeDropdown();
    };

    const moveActive = (direction) => {
      const options = visibleOptions.value;
      if (!options.length) return;
      const count = options.length;
      let nextIndex = activeIndex.value;
      if (nextIndex === -1) {
        nextIndex = direction > 0 ? 0 : count - 1;
      } else {
        nextIndex = (nextIndex + direction + count) % count;
      }
      setActiveIndex(nextIndex);
    };

    const handleTriggerKeydown = (event) => {
      switch (event.key) {
        case 'ArrowDown':
          event.preventDefault();
          if (!isOpen.value) {
            openDropdown();
          } else {
            moveActive(1);
          }
          break;
        case 'ArrowUp':
          event.preventDefault();
          if (!isOpen.value) {
            openDropdown();
            moveActive(-1);
          } else {
            moveActive(-1);
          }
          break;
        case 'Enter':
        case ' ': {
          event.preventDefault();
          if (!isOpen.value) {
            openDropdown();
          } else if (activeIndex.value >= 0) {
            const options = visibleOptions.value;
            if (activeIndex.value >= 0 && options[activeIndex.value]) {
              selectOption(options[activeIndex.value]);
            }
          }
          break;
        }
        case 'Escape':
          if (isOpen.value) {
            event.preventDefault();
            closeDropdown(false);
          }
          break;
        default:
      }
    };

    const handleListKeydown = (event) => {
      switch (event.key) {
        case 'ArrowDown':
          event.preventDefault();
          moveActive(1);
          break;
        case 'ArrowUp':
          event.preventDefault();
          moveActive(-1);
          break;
        case 'Home':
          event.preventDefault();
          setActiveIndex(0);
          break;
        case 'End':
          event.preventDefault();
          setActiveIndex(visibleOptions.value.length - 1);
          break;
        case 'Enter':
        case ' ': {
          event.preventDefault();
          const options = visibleOptions.value;
          if (activeIndex.value >= 0 && options[activeIndex.value]) {
            selectOption(options[activeIndex.value]);
          }
          break;
        }
        case 'Escape':
          event.preventDefault();
          closeDropdown();
          break;
        case 'Tab':
          closeDropdown(false);
          break;
        default:
      }
    };

    const handleClickOutside = (event) => {
      if (!isOpen.value) return;
      const target = event.target;
      const triggerEl = triggerRef.value;
      const listEl = listRef.value;
      if (triggerEl && triggerEl.contains(target)) return;
      if (listEl && listEl.contains(target)) return;
      closeDropdown(false);
    };

    onMounted(() => {
      document.addEventListener('click', handleClickOutside);
    });

    onBeforeUnmount(() => {
      document.removeEventListener('click', handleClickOutside);
    });

    const handleSearchKeydown = (event) => {
      if (event.key === 'ArrowDown') {
        event.preventDefault();
        if (visibleOptions.value.length) {
          setActiveIndex(0);
          listRef.value?.focus({ preventScroll: true });
        }
      } else if (event.key === 'ArrowUp') {
        event.preventDefault();
        if (visibleOptions.value.length) {
          setActiveIndex(visibleOptions.value.length - 1);
          listRef.value?.focus({ preventScroll: true });
        }
      } else if (event.key === 'Enter') {
        event.preventDefault();
        const options = visibleOptions.value;
        if (activeIndex.value >= 0 && options[activeIndex.value]) {
          selectOption(options[activeIndex.value]);
        }
      } else if (event.key === 'Escape') {
        if (searchQuery.value) {
          event.preventDefault();
          searchQuery.value = '';
        } else {
          closeDropdown();
        }
      }
    };

    const handleManageProjects = () => {
      emit('manage-projects');
      closeDropdown();
    };

    watch(visibleOptions, (options) => {
      if (!options.length) {
        activeIndex.value = -1;
        return;
      }
      if (activeIndex.value >= options.length) {
        activeIndex.value = options.length - 1;
        return;
      }
      if (activeIndex.value === -1) {
        if (selectedOption.value) {
          setActiveByValue(selectedOption.value.value);
          if (activeIndex.value === -1) {
            setActiveIndex(0);
          }
        } else {
          setActiveIndex(0);
        }
      }
    });

    watch(searchQuery, () => {
      nextTick(() => {
        if (!visibleOptions.value.length) {
          activeIndex.value = -1;
        } else if (activeIndex.value === -1) {
          setActiveIndex(0);
        }
      });
    });

    return () => {
      const rootClasses = ['hc-select'];
      if (isOpen.value) rootClasses.push('is-open');
      if (props.disabled) rootClasses.push('is-disabled');

      const triggerChildren = [];
      if (selectedOption.value) {
        triggerChildren.push(
          h('div', { class: 'hc-select__value' }, [
            h('span', { class: 'hc-select__label' }, selectedOption.value.label),
            selectedOption.value.description
              ? h('span', { class: 'hc-select__description' }, selectedOption.value.description)
              : null,
          ]),
        );
      } else {
        triggerChildren.push(
          h('div', { class: 'hc-select__value' }, [
            h('span', { class: 'hc-select__placeholder' }, props.placeholder),
          ]),
        );
      }

      triggerChildren.push(
        h(
          'span',
          { class: 'hc-select__icon', 'aria-hidden': 'true' },
          [
            h(
              'svg',
              {
                xmlns: 'http://www.w3.org/2000/svg',
                viewBox: '0 0 20 20',
                fill: 'none',
                width: 18,
                height: 18,
                stroke: 'currentColor',
                'stroke-width': 1.5,
              },
              [
                h('path', {
                  d: 'M6 8l4 4 4-4',
                  stroke: 'currentColor',
                  'stroke-width': 1.8,
                  'stroke-linecap': 'round',
                  'stroke-linejoin': 'round',
                }),
              ],
            ),
          ],
        ),
      );

      const dropdown =
        isOpen.value && props.options.length
          ? h('div', { class: 'hc-select__dropdown' }, [
              shouldShowSearch.value
                ? h('div', { class: 'hc-select__search' }, [
                    h('input', {
                      ref: searchInputRef,
                      value: searchQuery.value,
                      class: 'hc-select__search-input',
                      type: 'search',
                      placeholder: 'Search…',
                      onInput: (event) => (searchQuery.value = event.target.value),
                      onKeydown: handleSearchKeydown,
                    }),
                  ])
                : null,
              h(
                'ul',
                {
                  ref: listRef,
                  id: listboxId.value,
                  class: 'hc-select__options',
                  role: 'listbox',
                  tabindex: -1,
                  'aria-activedescendant':
                    activeIndex.value >= 0 ? optionId(activeIndex.value) : undefined,
                  onKeydown: handleListKeydown,
                },
                visibleOptions.value.length
                  ? visibleOptions.value.map((option, index) =>
                      h(
                        'li',
                        {
                          id: optionId(index),
                          key: option.value ?? index,
                          role: 'option',
                          class: [
                            'hc-select__option',
                            index === activeIndex.value ? 'is-active' : null,
                            selectedOption.value && option.value === selectedOption.value.value
                              ? 'is-selected'
                              : null,
                          ],
                          'aria-selected': selectedOption.value
                            ? String(option.value === selectedOption.value.value)
                            : 'false',
                          onClick: () => selectOption(option),
                          onMouseenter: () => setActiveIndex(index),
                        },
                        [
                          h('div', { class: 'hc-select__option-content' }, [
                            h('div', { class: 'hc-select__option-text' }, [
                              h('span', { class: 'hc-select__option-label' }, option.label),
                              option.description
                                ? h('span', { class: 'hc-select__option-description' }, option.description)
                                : null,
                            ]),
                            option.meta
                              ? h('span', { class: 'hc-select__option-meta' }, option.meta)
                              : null,
                          ]),
                        ],
                      ),
                    )
                  : [
                      h(
                        'li',
                        {
                          class: 'hc-select__empty',
                          role: 'presentation',
                        },
                        'No results found',
                      ),
                    ],
              ),
              h(
                'button',
                {
                  type: 'button',
                  class: 'hc-select__manage',
                  onClick: handleManageProjects,
                },
                'Manage Projects',
              ),
            ])
          : null;

      return h('div', { class: rootClasses.join(' ') }, [
        h(
          'button',
          {
            ref: triggerRef,
            id: componentId.value,
            type: 'button',
            class: 'hc-select__trigger',
            disabled: props.disabled,
            'aria-haspopup': 'listbox',
            'aria-expanded': String(isOpen.value),
            'aria-controls': listboxId.value,
            onClick: toggleDropdown,
            onKeydown: handleTriggerKeydown,
          },
          triggerChildren,
        ),
        dropdown,
      ]);
    };
  },
});
