import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StarRating from '../StarRating.vue'

describe('StarRating', () => {
  it('renders 5 buttons', () => {
    const wrapper = mount(StarRating, { props: { rating: 3 } })
    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(5)
  })

  it('emits update:rating when clicked and not readonly', async () => {
    const wrapper = mount(StarRating, { props: { rating: 2, readonly: false } })
    await wrapper.findAll('button')[3].trigger('click')
    expect(wrapper.emitted('update:rating')).toBeTruthy()
    expect(wrapper.emitted('update:rating')[0]).toEqual([4])
  })

  it('does not emit when readonly', async () => {
    const wrapper = mount(StarRating, { props: { rating: 3, readonly: true } })
    await wrapper.findAll('button')[0].trigger('click')
    expect(wrapper.emitted('update:rating')).toBeFalsy()
  })

  it('disables buttons when readonly', () => {
    const wrapper = mount(StarRating, { props: { rating: 3, readonly: true } })
    const buttons = wrapper.findAll('button')
    buttons.forEach(btn => {
      expect(btn.attributes('disabled')).toBeDefined()
    })
  })

  it('applies sm size class', () => {
    const wrapper = mount(StarRating, { props: { rating: 1, size: 'sm' } })
    const svg = wrapper.find('svg')
    expect(svg.classes()).toContain('w-4')
    expect(svg.classes()).toContain('h-4')
  })
})
