import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FilterChips from '../FilterChips.vue'

const options = [
  { value: 'all', label: 'All' },
  { value: 'road', label: 'Road' },
  { value: 'mtb', label: 'MTB' },
]

describe('FilterChips', () => {
  it('renders label and all options', () => {
    const wrapper = mount(FilterChips, {
      props: { label: 'Type', options, modelValue: 'all' },
    })
    expect(wrapper.text()).toContain('Type:')
    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(3)
    expect(buttons[0].text()).toBe('All')
    expect(buttons[1].text()).toBe('Road')
    expect(buttons[2].text()).toBe('MTB')
  })

  it('emits update:modelValue on click', async () => {
    const wrapper = mount(FilterChips, {
      props: { label: 'Type', options, modelValue: 'all' },
    })
    await wrapper.findAll('button')[1].trigger('click')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual(['road'])
  })

  it('highlights active option', () => {
    const wrapper = mount(FilterChips, {
      props: { label: 'Type', options, modelValue: 'road' },
    })
    const buttons = wrapper.findAll('button')
    expect(buttons[1].classes()).toContain('bg-primary')
    expect(buttons[0].classes()).not.toContain('bg-primary')
  })
})
