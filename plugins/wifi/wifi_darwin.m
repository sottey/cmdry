#import <CoreWLAN/CoreWLAN.h>
#import <CoreLocation/CoreLocation.h>
#import <Foundation/Foundation.h>

static void cmdry_request_location_access(void) {
  if (@available(macOS 10.15, *)) {
    CLLocationManager *manager = [[CLLocationManager alloc] init];
    if (manager.authorizationStatus != kCLAuthorizationStatusNotDetermined) return;

    [manager requestWhenInUseAuthorization];
    NSDate *deadline = [NSDate dateWithTimeIntervalSinceNow:6.0];
    while (manager.authorizationStatus == kCLAuthorizationStatusNotDetermined && [deadline timeIntervalSinceNow] > 0) {
      [[NSRunLoop currentRunLoop] runMode:NSDefaultRunLoopMode beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
    }
  }
}

char *cmdry_wifi_info(void) {
  @autoreleasepool {
    cmdry_request_location_access();
    CWWiFiClient *client = [CWWiFiClient sharedWiFiClient];
    CWInterface *interface = [client interface];
    if (interface == nil) {
      for (CWInterface *candidate in client.interfaces) {
        if (candidate.ssid != nil) {
          interface = candidate;
          break;
        }
        if (interface == nil) {
          interface = candidate;
        }
      }
    }
    if (interface == nil) return NULL;
    NSString *ssid = interface.ssid;
    NSString *name = interface.interfaceName ?: @"";
    NSString *bssid = interface.bssid ?: @"";
    NSString *signal = [NSString stringWithFormat:@"%ld dBm", (long)interface.rssiValue];
    NSString *result = [NSString stringWithFormat:@"%@\t%@\t%@\t%@\t", name, ssid, bssid, signal];
    return strdup(result.UTF8String);
  }
}
