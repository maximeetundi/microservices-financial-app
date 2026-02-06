<template>
  <div class="space-y-8">
    <!-- Breadcrumb -->
    <nav class="flex items-center text-sm text-gray-500 dark:text-gray-400 overflow-x-auto whitespace-nowrap pb-2">
      <NuxtLink :to="`/shops/${shopSlug}`" class="hover:text-indigo-600 transition-colors">
        {{ shop?.name || 'Boutique' }}
      </NuxtLink>
      <span class="mx-2">›</span>
      <NuxtLink
        v-if="product?.category_name"
        :to="`/shops/${shopSlug}?category=${product.category_slug}`"
        class="hover:text-indigo-600 transition-colors"
      >
        {{ product.category_name }}
      </NuxtLink>
      <span v-if="product?.category_name" class="mx-2">›</span>
      <span class="text-gray-900 dark:text-white font-medium truncate max-w-[200px]">
        {{ product?.name || 'Produit' }}
      </span>
    </nav>


    <!-- Loading State -->
    <div v-if="loading" class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <div class="animate-pulse bg-gray-200 dark:bg-slate-800 aspect-square rounded-2xl"></div>
      <div class="space-y-4">
        <div class="animate-pulse bg-gray-200 dark:bg-slate-800 h-8 w-3/4 rounded-lg"></div>
        <div class="animate-pulse bg-gray-200 dark:bg-slate-800 h-6 w-1/3 rounded-lg"></div>
        <div class="animate-pulse bg-gray-200 dark:bg-slate-800 h-24 w-full rounded-lg"></div>
        <div class="animate-pulse bg-gray-200 dark:bg-slate-800 h-12 w-full rounded-lg"></div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error || !product" class="text-center py-20 bg-white dark:bg-slate-900 rounded-2xl">
      <div class="w-24 h-24 mx-auto mb-4 bg-gray-100 dark:bg-slate-800 rounded-full flex items-center justify-center text-5xl">😕</div>
      <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Produit introuvable</h2>
      <p class="text-gray-500 mb-6">Ce produit n'existe pas ou a été supprimé.</p>
      <NuxtLink
        :to="`/shops/${shopSlug}`"
        class="inline-block px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-xl transition-colors"
      >
        Retour à la boutique
      </NuxtLink>
    </div>

    <!-- Product Content -->
    <template v-else>
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Image Gallery -->
        <div class="space-y-4">
          <!-- Main Image -->
          <div class="relative aspect-square bg-gray-100 dark:bg-slate-800 rounded-2xl overflow-hidden group">
            <img
              v-if="product.images?.length"
              :src="product.images[selectedImageIndex]"
              :alt="product.name"
              class="w-full h-full object-contain p-4 group-hover:scale-105 transition-transform duration-500"
            >
            <div v-else class="w-full h-full flex items-center justify-center text-7xl text-gray-300">📦</div>

            <!-- Zoom Button -->
            <button
              v-if="product.images?.length"
              @click="showLightbox = true"
              class="absolute top-4 right-4 w-10 h-10 bg-white/90 dark:bg-slate-800/90 rounded-xl shadow-sm flex items-center justify-center text-gray-600 dark:text-gray-400 hover:bg-white dark:hover:bg-slate-700 transition-colors opacity-0 group-hover:opacity-100"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7"/>
              </svg>
            </button>

            <!-- Badges -->
            <div class="absolute top-4 left-4 flex flex-col gap-2">
              <span v-if="discountPercent > 0" class="px-3 py-1.5 bg-red-500 text-white text-xs font-bold rounded-lg shadow">
                -{{ discountPercent }}%
              </span>
              <span v-if="product.is_featured" class="px-3 py-1.5 bg-amber-500 text-white text-xs font-bold rounded-lg shadow">
                ⭐ Vedette
              </span>
              <span v-if="product.stock === 0" class="px-3 py-1.5 bg-gray-800 text-white text-xs font-bold rounded-lg shadow">
                Épuisé
              </span>
            </div>

            <!-- Navigation Arrows -->
            <template v-if="product.images?.length > 1">
              <button
                @click="prevImage"
                class="absolute left-4 top-1/2 -translate-y-1/2 w-10 h-10 bg-white/90 dark:bg-slate-800/90 rounded-full shadow flex items-center justify-center text-gray-600 hover:bg-white transition-colors opacity-0 group-hover:opacity-100"
              >
                ‹
              </button>
              <button
                @click="nextImage"
                class="absolute right-4 top-1/2 -translate-y-1/2 w-10 h-10 bg-white/90 dark:bg-slate-800/90 rounded-full shadow flex items-center justify-center text-gray-600 hover:bg-white transition-colors opacity-0 group-hover:opacity-100"
              >
                ›
              </button>
            </template>
          </div>

          <!-- Thumbnails -->
          <div v-if="product.images?.length > 1" class="flex gap-3 overflow-x-auto pb-2 scrollbar-hide">
            <button
              v-for="(img, i) in product.images"
              :key="i"
              @click="selectedImageIndex = i"
              class="flex-shrink-0 w-16 h-16 rounded-xl overflow-hidden border-2 transition-all"
              :class="selectedImageIndex === i ? 'border-indigo-500 shadow-lg scale-105' : 'border-transparent opacity-70 hover:opacity-100'"
            >
              <img :src="img" class="w-full h-full object-cover" :alt="`Image ${i + 1}`">
            </button>
          </div>
        </div>

        <!-- Product Info -->
        <div class="space-y-6">
          <!-- Category & Stock -->
          <div class="flex items-center gap-3">
            <span v-if="product.category_name" class="text-xs font-medium text-indigo-600 dark:text-indigo-400 uppercase tracking-wider">
              {{ product.category_name }}
            </span>
            <span
              v-if="product.stock > 0"
              class="px-2.5 py-1 bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 text-xs font-medium rounded-full"
            >
              ✓ En stock ({{ product.stock }})
            </span>
            <span v-else class="px-2.5 py-1 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 text-xs font-medium rounded-full">
              Rupture de stock
            </span>
          </div>

          <!-- Name -->
          <h1 class="text-2xl md:text-3xl font-bold text-gray-900 dark:text-white leading-tight">
            {{ product.name }}
          </h1>

          <!-- Rating Summary -->
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <div class="flex text-amber-400 text-lg">
                <span v-for="i in 5" :key="i">{{ i <= Math.round(averageRating) ? '★' : '☆' }}</span>
              </div>
              <span class="text-lg font-bold text-gray-900 dark:text-white">{{ averageRating.toFixed(1) }}</span>
            </div>
            <button @click="scrollToReviews" class="text-sm text-indigo-600 dark:text-indigo-400 hover:underline">
              {{ reviews.length }} avis
            </button>
          </div>

          <!-- Price -->
          <div class="flex items-end gap-4">
            <span class="text-3xl md:text-4xl font-bold text-indigo-600 dark:text-indigo-400">
              {{ formatPrice(product.price) }}
            </span>
            <span
              v-if="product.compare_at_price && product.compare_at_price > product.price"
              class="text-xl text-gray-400 line-through mb-1"
            >
              {{ formatPrice(product.compare_at_price) }}
            </span>
            <span
              v-if="discountPercent > 0"
              class="px-3 py-1 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 font-bold text-sm rounded-lg mb-1"
            >
              Économisez {{ formatPrice(product.compare_at_price - product.price) }}
            </span>
          </div>

          <!-- Description -->
          <p class="text-gray-600 dark:text-gray-300 leading-relaxed">
            {{ product.description }}
          </p>

          <!-- Tags -->
          <div v-if="product.tags && product.tags.length > 0" class="flex flex-wrap gap-2 mt-4">
            <span 
              v-for="tag in product.tags" 
              :key="tag"
              class="px-2.5 py-1 bg-gray-100 dark:bg-slate-700 text-gray-600 dark:text-gray-300 text-xs rounded-md font-medium"
            >
              #{{ tag }}
            </span>
          </div>

          <!-- Quantity Selector -->
          <div class="flex items-center gap-4 p-4 bg-gray-50 dark:bg-slate-800 rounded-xl">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Quantité</span>
            <div class="flex items-center gap-2">
              <button
                @click="quantity > 1 && quantity--"
                class="w-10 h-10 rounded-lg bg-white dark:bg-slate-700 border border-gray-200 dark:border-gray-600 flex items-center justify-center text-lg font-bold hover:border-indigo-500 transition-colors"
              >
                −
              </button>
              <input
                v-model.number="quantity"
                type="number"
                min="1"
                :max="product.stock || 99"
                class="w-16 h-10 text-center border border-gray-200 dark:border-gray-600 rounded-lg bg-white dark:bg-slate-700 font-bold"
              >
              <button
                @click="quantity < (product.stock || 99) && quantity++"
                class="w-10 h-10 rounded-lg bg-white dark:bg-slate-700 border border-gray-200 dark:border-gray-600 flex items-center justify-center text-lg font-bold hover:border-indigo-500 transition-colors"
              >
                +
              </button>
            </div>
            <span v-if="product.stock && product.stock < 10" class="text-xs text-amber-600 dark:text-amber-400 ml-auto">
              ⚠️ Plus que {{ product.stock }} en stock
            </span>
          </div>

          <!-- Action Buttons -->
          <div class="flex flex-col sm:flex-row gap-3">
            <button
              @click="addToCart"
              :disabled="product.stock === 0 || addingToCart"
              class="flex-1 py-4 px-6 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 shadow-lg shadow-indigo-500/25"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/>
              </svg>
              {{ addingToCart ? 'Ajout...' : product.stock === 0 ? 'Rupture de stock' : 'Ajouter au panier' }}
            </button>
            <button
              @click="toggleFavorite"
              class="p-4 rounded-xl border-2 transition-all"
              :class="isFavorite ? 'border-red-500 bg-red-50 dark:bg-red-900/20 text-red-500' : 'border-gray-200 dark:border-gray-700 text-gray-400 hover:border-red-300 hover:text-red-500'"
            >
              <svg v-if="isFavorite" class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
              </svg>
              <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
              </svg>
            </button>
            <button
              @click="shareProduct"
              class="p-4 rounded-xl border-2 border-gray-200 dark:border-gray-700 text-gray-400 hover:border-indigo-300 hover:text-indigo-500 transition-all"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
              </svg>
            </button>
          </div>

          <button
            type="button"
            @click="contactSeller"
            class="w-full py-3 px-6 bg-white dark:bg-slate-900 border border-gray-200 dark:border-gray-700 rounded-xl font-semibold text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors"
          >
            Contacter le vendeur
          </button>

          <!-- Trust Badges -->
          <div class="grid grid-cols-2 gap-3 pt-4 border-t border-gray-100 dark:border-gray-800">
            <div
              v-for="badge in visibleTrustBadges"
              :key="badge.key"
              class="flex items-center gap-3 p-3 bg-gray-50 dark:bg-slate-800 rounded-xl"
            >
              <span class="text-2xl">{{ badge.icon }}</span>
              <div>
                <p class="text-xs font-semibold text-gray-900 dark:text-white">{{ badge.title }}</p>
                <p class="text-[10px] text-gray-500">{{ badge.subtitle }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Tabs Section -->
      <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-gray-800 overflow-hidden">
        <!-- Tab Headers -->
        <div class="flex border-b border-gray-100 dark:border-gray-800">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            class="flex-1 py-4 px-6 text-sm font-semibold transition-colors relative"
            :class="activeTab === tab.id ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'"
          >
            {{ tab.label }}
            <span v-if="tab.count" class="ml-1 text-xs">({{ tab.count }})</span>
            <div v-if="activeTab === tab.id" class="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600"></div>
          </button>
        </div>

        <!-- Tab Content -->
        <div class="p-6">
          <!-- Description Tab -->
          <div v-if="activeTab === 'description'" class="prose dark:prose-invert max-w-none">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-4">À propos de ce produit</h3>
            <p class="text-gray-600 dark:text-gray-300 leading-relaxed whitespace-pre-line">{{ product.description }}</p>

            <!-- Product Details -->
            <div v-if="product.sku || product.weight || product.dimensions" class="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4">
              <div v-if="product.sku" class="p-4 bg-gray-50 dark:bg-slate-800 rounded-xl">
                <p class="text-xs text-gray-500 mb-1">SKU</p>
                <p class="font-medium text-gray-900 dark:text-white">{{ product.sku }}</p>
              </div>
              <div v-if="product.weight" class="p-4 bg-gray-50 dark:bg-slate-800 rounded-xl">
                <p class="text-xs text-gray-500 mb-1">Poids</p>
                <p class="font-medium text-gray-900 dark:text-white">{{ product.weight }}</p>
              </div>
            </div>
          </div>

          <!-- Reviews Tab -->
          <div v-if="activeTab === 'reviews'" id="reviews-section" class="space-y-6">
            <!-- Rating Overview -->
            <div class="flex flex-col md:flex-row gap-6 p-6 bg-gray-50 dark:bg-slate-800 rounded-xl">
              <div class="text-center">
                <div class="text-5xl font-bold text-gray-900 dark:text-white mb-2">{{ averageRating.toFixed(1) }}</div>
                <div class="flex justify-center text-amber-400 text-xl mb-1">
                  <span v-for="i in 5" :key="i">{{ i <= Math.round(averageRating) ? '★' : '☆' }}</span>
                </div>
                <p class="text-sm text-gray-500">{{ reviews.length }} avis</p>
              </div>
              <div class="flex-1 space-y-2">
                <div v-for="rating in [5, 4, 3, 2, 1]" :key="rating" class="flex items-center gap-3">
                  <span class="text-sm text-gray-600 dark:text-gray-400 w-8">{{ rating }}★</span>
                  <div class="flex-1 h-2 bg-gray-200 dark:bg-slate-700 rounded-full overflow-hidden">
                    <div
                      class="h-full bg-amber-400 rounded-full transition-all"
                      :style="{ width: `${getRatingPercent(rating)}%` }"
                    ></div>
                  </div>
                  <span class="text-sm text-gray-500 w-10 text-right">{{ getRatingCount(rating) }}</span>
                </div>
              </div>
            </div>

            <!-- Write Review -->
            <div class="p-6 border border-gray-200 dark:border-gray-700 rounded-xl">
              <h4 class="font-bold text-gray-900 dark:text-white mb-4">Donnez votre avis</h4>
              <div class="flex items-center gap-1 mb-4">
                <button
                  v-for="i in 5"
                  :key="i"
                  @click="newReview.rating = i"
                  class="text-3xl transition-transform hover:scale-110"
                  :class="i <= newReview.rating ? 'text-amber-400' : 'text-gray-300 dark:text-gray-600'"
                >
                  ★
                </button>
                <span class="ml-3 text-sm text-gray-500">
                  {{ newReview.rating > 0 ? ratingLabels[newReview.rating - 1] : 'Cliquez pour noter' }}
                </span>
              </div>
              <textarea
                v-model="newReview.comment"
                rows="3"
                class="w-full px-4 py-3 border border-gray-200 dark:border-gray-700 rounded-xl bg-white dark:bg-slate-800 text-gray-900 dark:text-white resize-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Partagez votre expérience avec ce produit..."
              ></textarea>
              <div class="flex justify-end mt-3">
                <button
                  @click="submitReview"
                  :disabled="!newReview.rating || submittingReview"
                  class="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-xl transition-colors disabled:opacity-50"
                >
                  {{ submittingReview ? 'Envoi...' : 'Publier mon avis' }}
                </button>
              </div>
            </div>

            <!-- Reviews List -->
            <div class="space-y-4">
              <div
                v-for="review in reviews"
                :key="review.id"
                class="p-5 border border-gray-100 dark:border-gray-800 rounded-xl"
              >
                <div class="flex items-start gap-4">
                  <div class="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center text-white font-bold flex-shrink-0">
                    {{ review.user_name?.charAt(0)?.toUpperCase() || 'U' }}
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-3 mb-1">
                      <span class="font-semibold text-gray-900 dark:text-white">{{ review.user_name || 'Utilisateur' }}</span>
                      <span v-if="review.verified" class="px-2 py-0.5 bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400 text-[10px] font-bold rounded-full">
                        ✓ Achat vérifié
                      </span>
                    </div>
                    <div class="flex items-center gap-2 mb-3">
                      <div class="flex text-amber-400 text-sm">
                        <span v-for="i in 5" :key="i">{{ i <= review.rating ? '★' : '☆' }}</span>
                      </div>
                      <span class="text-xs text-gray-400">{{ formatDate(review.created_at) }}</span>
                    </div>
                    <p class="text-gray-600 dark:text-gray-300 text-sm leading-relaxed">{{ review.comment }}</p>

                    <!-- Review Actions -->
                    <div class="flex items-center gap-4 mt-3">
                      <button
                        @click="likeReview(review.id)"
                        class="flex items-center gap-1.5 text-xs text-gray-500 hover:text-indigo-600 transition-colors"
                      >
                        <span>👍</span> Utile ({{ review.likes || 0 }})
                      </button>
                      <button class="text-xs text-gray-500 hover:text-indigo-600 transition-colors">
                        Signaler
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Empty Reviews -->
              <div v-if="reviews.length === 0" class="text-center py-12">
                <div class="w-16 h-16 mx-auto mb-4 bg-gray-100 dark:bg-slate-800 rounded-full flex items-center justify-center text-3xl">💬</div>
                <h4 class="font-bold text-gray-900 dark:text-white mb-2">Aucun avis pour le moment</h4>
                <p class="text-sm text-gray-500">Soyez le premier à donner votre avis sur ce produit !</p>
              </div>
            </div>
          </div>

          <!-- Shipping Tab -->
          <div v-if="activeTab === 'shipping'" class="space-y-4">
            <div class="flex items-start gap-4 p-4 bg-gray-50 dark:bg-slate-800 rounded-xl">
              <span class="text-2xl">📦</span>
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-1">Livraison Standard</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">Livraison en 3-5 jours ouvrés. Gratuite à partir de 25 000 FCFA.</p>
              </div>
            </div>
            <div class="flex items-start gap-4 p-4 bg-gray-50 dark:bg-slate-800 rounded-xl">
              <span class="text-2xl">⚡</span>
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-1">Livraison Express</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">Livraison en 24-48h. Frais de 2 500 FCFA.</p>
              </div>
            </div>
            <div class="flex items-start gap-4 p-4 bg-gray-50 dark:bg-slate-800 rounded-xl">
              <span class="text-2xl">🏪</span>
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-1">Retrait en point relais</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">Récupérez votre commande dans l'un de nos points relais partenaires.</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Similar Products -->
      <section v-if="similarProducts.length > 0">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">Produits similaires</h2>
          <NuxtLink
            :to="`/shops/${shopSlug}?category=${product.category_slug}`"
            class="text-sm text-indigo-600 dark:text-indigo-400 font-medium hover:underline"
          >
            Voir plus →
          </NuxtLink>
        </div>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <ProductCard
            v-for="product in similarProducts.slice(0, 4)"
            :key="product.id"
            :product="product"
            :shop-slug="shopSlug"
            @add-to-cart="addToCart"
          />
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { inject, ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi, type Shop, type Product, type Review, type ShopTrustBadge } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'
import { messagingAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const cartStore = useCartStore()
const authStore = useAuthStore()

// ...

const contactSeller = async () => {
  await authStore.initializeAuth()
  if (!authStore.isAuthenticated || !authStore.user?.id) {
    router.push(`/auth/login?redirect=${encodeURIComponent(route.fullPath)}`)
    return
  }

  const ownerId = shop.value?.owner_id
  if (!ownerId) return

  try {
    const myName = `${authStore.user.firstName || ''} ${authStore.user.lastName || ''}`.trim() || authStore.user.email
    const res = await messagingAPI.createConversation({
      participant_id: ownerId,
      participant_name: shop.value?.name || 'Vendeur',
      my_name: myName,
      context: {
        type: 'shop',
        shop_id: shop.value?.id,
        shop_name: shop.value?.name,
        product_id: product.value?.id,
        product_name: product.value?.name,
      }
    })
    const convId = res.data?.id
    if (convId) {
      router.push(`/messages?conversation_id=${encodeURIComponent(convId)}`)
    } else {
      router.push('/messages')
    }
  } catch (e) {
    console.error('Failed to create conversation with seller', e)
    router.push('/messages')
  }
}

// ...

const loadData = async () => {
  loading.value = true
  // ...
  
  try {
    product.value = await shopApi.getProduct(shopSlug.value, productSlug.value)
    
    // Load reviews
    const reviewsData = await shopApi.listReviews(product.value.id)
    reviews.value = reviewsData.reviews || []
    
    // Load similar products
    const similar = await shopApi.listProducts(shopSlug.value, 1, 4, { category: product.value.category_slug })
    similarProducts.value = similar.products?.filter((p: Product) => p.id !== product.value?.id).slice(0, 4) || []
    
  } catch (e) {
    console.error('Failed to load product', e)
    error.value = true
  } finally {
    loading.value = false
  }
}

const submitReview = async () => {
  if (!product.value) return
  
  submittingReview.value = true
  try {
    const review = await shopApi.createReview(product.value.id, {
      rating: newReview.value.rating,
      comment: newReview.value.comment
    })
    
    reviews.value.unshift(review)
    newReview.value = { rating: 0, comment: '' }
    
    // Update local rating
    // In real app, we might want to reload product to get updated average
  } catch (e) {
    console.error('Failed to submit review', e)
    alert('Erreur lors de l\'envoi de l\'avis')
  } finally {
    submittingReview.value = false
  }
}

watch([shopSlug, productSlug], loadData, { immediate: true })
</script>
