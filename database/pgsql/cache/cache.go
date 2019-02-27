package cache

import lru "github.com/hashicorp/golang-lru"

// FeatureCache caches the feature as key and feature IDs as value.
type FeatureCache struct {
	c *lru.Cache
}

// NamespaceCache caches the namespace as key and namespace IDs as value.
type NamespaceCache struct {
	c *lru.Cache
}

// NamespacedFeatureCache caches the namespaced feature as key and
// namespacedFeature ID as value.
type NamespacedFeatureCache struct {
	c *lru.Cache
}

// PgCache caches the database immutable entities.
type PgCache struct {
	features          FeatureCache
	namespaces        NamespaceCache
	namespacedFeature NamespacedFeatureCache
}
