#include "fat/socket.h"

#include "fatstd_go.h"

fat_Status fat_TcpDialUTF8(const char *addr, fat_TcpConn *out_conn, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tcp_dial_utf8((char *)addr, (uintptr_t *)out_conn, (uintptr_t *)out_err);
}

fat_Status fat_TcpListenerListenUTF8(const char *addr, fat_TcpListener *out_listener, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tcp_listener_listen_utf8((char *)addr, (uintptr_t *)out_listener, (uintptr_t *)out_err);
}

fat_Status fat_TcpListenerAccept(fat_TcpListener listener, fat_TcpConn *out_conn, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tcp_listener_accept((uintptr_t)listener, (uintptr_t *)out_conn, (uintptr_t *)out_err);
}

fat_String fat_TcpListenerAddr(fat_TcpListener listener) {
  return (fat_String)fatstd_go_tcp_listener_addr((uintptr_t)listener);
}

fat_Status fat_TcpListenerClose(fat_TcpListener listener, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tcp_listener_close((uintptr_t)listener, (uintptr_t *)out_err);
}

fat_Status fat_TcpConnRead(
  fat_TcpConn conn,
  void *dst,
  size_t dst_len,
  size_t *out_n,
  bool *out_eof,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_tcp_conn_read((uintptr_t)conn, (char *)dst, dst_len, out_n, (_Bool *)out_eof,
                                             (uintptr_t *)out_err);
}

fat_Status fat_TcpConnWrite(fat_TcpConn conn, const void *src, size_t src_len, size_t *out_n, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tcp_conn_write((uintptr_t)conn, (char *)src, src_len, out_n, (uintptr_t *)out_err);
}

fat_String fat_TcpConnLocalAddr(fat_TcpConn conn) {
  return (fat_String)fatstd_go_tcp_conn_local_addr((uintptr_t)conn);
}

fat_String fat_TcpConnRemoteAddr(fat_TcpConn conn) {
  return (fat_String)fatstd_go_tcp_conn_remote_addr((uintptr_t)conn);
}

fat_Status fat_TcpConnClose(fat_TcpConn conn, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tcp_conn_close((uintptr_t)conn, (uintptr_t *)out_err);
}

fat_Status fat_UdpListenUTF8(const char *addr, fat_UdpConn *out_conn, fat_Error *out_err) {
  return (fat_Status)fatstd_go_udp_listen_utf8((char *)addr, (uintptr_t *)out_conn, (uintptr_t *)out_err);
}

fat_Status fat_UdpDialUTF8(const char *addr, fat_UdpConn *out_conn, fat_Error *out_err) {
  return (fat_Status)fatstd_go_udp_dial_utf8((char *)addr, (uintptr_t *)out_conn, (uintptr_t *)out_err);
}

fat_Status fat_UdpConnReadFrom(
  fat_UdpConn conn,
  void *dst,
  size_t dst_len,
  size_t *out_n,
  fat_String *out_addr,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_udp_conn_read_from((uintptr_t)conn, (char *)dst, dst_len, out_n,
                                                  (uintptr_t *)out_addr, (uintptr_t *)out_err);
}

fat_Status fat_UdpConnWriteToUTF8(
  fat_UdpConn conn,
  const void *src,
  size_t src_len,
  const char *addr,
  size_t *out_n,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_udp_conn_write_to_utf8((uintptr_t)conn, (char *)src, src_len, (char *)addr, out_n,
                                                     (uintptr_t *)out_err);
}

fat_Status fat_UdpConnWrite(fat_UdpConn conn, const void *src, size_t src_len, size_t *out_n, fat_Error *out_err) {
  return (fat_Status)fatstd_go_udp_conn_write((uintptr_t)conn, (char *)src, src_len, out_n, (uintptr_t *)out_err);
}

fat_String fat_UdpConnLocalAddr(fat_UdpConn conn) {
  return (fat_String)fatstd_go_udp_conn_local_addr((uintptr_t)conn);
}

fat_String fat_UdpConnRemoteAddr(fat_UdpConn conn) {
  return (fat_String)fatstd_go_udp_conn_remote_addr((uintptr_t)conn);
}

fat_Status fat_UdpConnClose(fat_UdpConn conn, fat_Error *out_err) {
  return (fat_Status)fatstd_go_udp_conn_close((uintptr_t)conn, (uintptr_t *)out_err);
}
