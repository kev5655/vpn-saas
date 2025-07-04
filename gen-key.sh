wg genkey | tee client_priv.key | wg pubkey > client_pub.key
export CLIENT_PUBKEY="$(<client_pub.key)"
