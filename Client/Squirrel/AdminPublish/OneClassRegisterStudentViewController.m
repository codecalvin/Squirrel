//
//  OneClassRegisterStudentViewController.m
//  Squirrel
//
//  Created by JamesMao on 6/29/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "OneClassRegisterStudentViewController.h"
#import "Define.h"
#import "DataModel/OneClassReigsterStudentData.h"
#import "AdminPublishTableViewCell.h"
#import "DataModel/UserRegisterDataItem.h"

@interface OneClassRegisterStudentViewController ()
{
    IBOutlet UITableView* studentTableView_;
    NSString* classUniqueID_;
}
@end

@implementation OneClassRegisterStudentViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view from its nib.
}

- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}


- (void)viewWillAppear:(BOOL)animated
{
    NSString* urlString = [NSString stringWithFormat:@"%@%@%@", SERVER_IP, URL_PART_ONE_CLASS_USERS, [OneClassReigsterStudentData singleton].classUniqueID];
    [self request:RequestTypeGet urlString:urlString parameters:nil];
    [[OneClassReigsterStudentData singleton] reset];
}

- (void)onSuccess:(AFHTTPRequestOperation *)operation responseObject:(id)responseObject
{
    [super onSuccess:operation responseObject:responseObject];
    
    [[OneClassReigsterStudentData singleton] reset];
    [[OneClassReigsterStudentData singleton] setJsonDictionaryData:responseObject];
    [studentTableView_ reloadData];
}

- (void)onFail:(AFHTTPRequestOperation *)operation error:(NSError*)error
{
    [super onFail:operation error:error];
    
    [[OneClassReigsterStudentData singleton] reset];
    [studentTableView_ reloadData];
}


#pragma mark -
#pragma mark Table View Data Source Methods
-(NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section
{
    int count = [[OneClassReigsterStudentData singleton] getUserRegisterDataItemCount];
    return count;
}

-(UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath
{
    static NSString *cellIdentifier = @"cellIdentifier";
    
    UITableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:cellIdentifier];
    if (cell == nil)
    {
        cell = [self getTableViewCellbyLoadingNib];
    }
    
    NSUInteger row = [indexPath row];
    
    NSMutableDictionary* dataItemEelements =  [[OneClassReigsterStudentData singleton] getUserRegisterDataItem:(int)row];
    AdminPublishTableViewCell* adminPublishTableViewCell = (AdminPublishTableViewCell*)cell;
    
    NSString* userName = [UserRegisterDataItem getUserName:dataItemEelements];
    if (userName == nil)
    {
        userName = @"";
    }
    
    [adminPublishTableViewCell setClassName:userName];
    [adminPublishTableViewCell setClassTime:@""];
    return cell;
}

- (UITableViewCell*)getTableViewCellbyLoadingNib
{
    NSString* personTableViewCell = @"AdminPublishTableViewCell";
    
    AdminPublishTableViewCell *cell;
    NSArray *nib = [[NSBundle mainBundle] loadNibNamed:personTableViewCell
                                                 owner:self
                                               options:nil];
    for (id oneObject in nib)
    {
        if ([oneObject isKindOfClass:[AdminPublishTableViewCell class]])
        {
            cell = (AdminPublishTableViewCell *)oneObject;
        }
    }
    return cell;
}

-(CGFloat)tableView:(UITableView*)tableView heightForRowAtIndexPath:(NSIndexPath*)indexPath
{
    return [self getTableViewCellbyLoadingNib].frame.size.height;
}

#pragma mark -
#pragma mark Table Delegate Methods
-(void)tableView:(UITableView*)tableView didSelectRowAtIndexPath:(NSIndexPath*)indexPath
{
    // to do 
}


@end
