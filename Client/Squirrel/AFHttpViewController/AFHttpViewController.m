//
//  AFHttpViewController.m
//  Squirrel
//
//  Created by JamesMao on 7/5/15.
//
//

#import "AFHttpViewController.h"
#import "AFNetworking.h"

@interface AFHttpViewController ()

@end

@implementation AFHttpViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    
}

- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
}

- (void)request:(RequestType)type urlString:(NSString*)urlString parameters:(id)parameters
{
    AFHTTPRequestOperationManager *manager = [AFHTTPRequestOperationManager manager];
    
    [self setHttpsSecurityPolicy:urlString operationManager:manager];
    
    switch (type) {
        case RequestTypeGet:
        {
            NSLog(@"request GET:%@", urlString);
            [manager GET:urlString
               parameters:parameters
                  success:^(AFHTTPRequestOperation *operation, id responseObject) {
                      [self onSuccess:operation responseObject:responseObject];}
                  failure:^(AFHTTPRequestOperation *operation, NSError *error) {
                      [self onFail:operation error:error];}
             ];
            break;
        }
        case RequestTypePost:
        {
            NSLog(@"request POST:%@", urlString);
            [manager POST:urlString
               parameters:parameters
                  success:^(AFHTTPRequestOperation *operation, id responseObject) {
                      [self onSuccess:operation responseObject:responseObject];}
                  failure:^(AFHTTPRequestOperation *operation, NSError *error) {
                      [self onFail:operation error:error];}
             ];
            break;
        }
        case RequestTypeDelete:
        {
            NSLog(@"request DELETE:%@", urlString);
            [manager DELETE:urlString
               parameters:parameters
                  success:^(AFHTTPRequestOperation *operation, id responseObject) {
                      [self onSuccess:operation responseObject:responseObject];}
                  failure:^(AFHTTPRequestOperation *operation, NSError *error) {
                      [self onFail:operation error:error];}
             ];
            break;        }
        default:
            break;
    }
}

- (void)setHttpsSecurityPolicy:(NSString*)urlString operationManager:(AFHTTPRequestOperationManager*)manager
{
    if ([urlString rangeOfString:@"https://"].length > 0)
    {
        AFSecurityPolicy *securityPolicy = [[AFSecurityPolicy alloc] init];
        [securityPolicy setAllowInvalidCertificates:YES];
        [manager setSecurityPolicy:securityPolicy];
    }
}


- (void)onSuccess:(AFHTTPRequestOperation *)operation responseObject:(id)responseObject
{
    NSLog(@"success url: %@",[[[operation response] URL] absoluteString]);
    NSLog(@"success data: %@", responseObject);
    
    // to be override
}

- (void)onFail:(AFHTTPRequestOperation *)operation error:(NSError*)error
{
    NSLog(@"Error url: %@",[[[operation response] URL] absoluteString]);
     NSLog(@"Error: %@", error);
    
    // to be override
}

@end
