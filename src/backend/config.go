// config.go

package backend

const ENABLED  = true
const DISABLED = false

/**
 * ENABLE_SUBSCRIPTION_EXECUTOR
 * 
 * This option enables an internal subscription operation executor. The internal 
 * subscription implements just for small applications only. For supporting 
 * subscriptions at scale, we suggest implement subscription by yourself or
 * use a dedicated subscription/push platform.  
 *
 * IMPORTANT: DISABLE this option can improve performance.
 */
const ENABLE_SUBSCRIPTION_EXECUTOR = ENABLED



/**
 * [ENABLE_JIT description]
 * @type {[type]}
 */
const ENABLE_JIT = ENABLED