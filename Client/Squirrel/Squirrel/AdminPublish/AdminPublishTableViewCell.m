//
//  AdminPublishTableViewCell.m
//  Squirrel
//
//  Created by JamesMao on 4/12/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "AdminPublishTableViewCell.h"

@interface AdminPublishTableViewCell()
{
    IBOutlet UILabel* className_;
    IBOutlet UILabel* classTime_;
}

@end

@implementation AdminPublishTableViewCell

- (void)awakeFromNib
{
    // Initialization code
}

- (void)setSelected:(BOOL)selected animated:(BOOL)animated
{
    [super setSelected:selected animated:animated];

    // Configure the view for the selected state
}


- (void)setClassName:(NSString*)className
{
    className_.text = className;
}

- (void)setClassTime:(NSString*)classTime
{
    classTime_.text = classTime;
}

@end
